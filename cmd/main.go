package main

import (
	"fmt"
	"github.com/lazylex/watch-store-store/internal/adapters/message_broker/kafka"
	restServer "github.com/lazylex/watch-store-store/internal/adapters/rest/server"
	"github.com/lazylex/watch-store-store/internal/config"
	"github.com/lazylex/watch-store-store/internal/logger"
	prometheusMetrics "github.com/lazylex/watch-store-store/internal/metrics"
	"github.com/lazylex/watch-store-store/internal/ports/repository"
	"github.com/lazylex/watch-store-store/internal/repository/mysql"
	"github.com/lazylex/watch-store-store/internal/service"
	"github.com/lazylex/watch-store-store/pkg/db-viewer/db_reader"
	"github.com/lazylex/watch-store-store/pkg/secure"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	cfg := config.MustLoad()
	slog.SetDefault(logger.MustCreate(cfg.Env, cfg.Instance))
	if err := clearScreen(); err != nil {
		slog.Error(err.Error())
	}

	var permissionsChan chan secure.NameNumber
	if cfg.Env != config.EnvironmentLocal {
		permissionsChan = make(chan secure.NameNumber)
		appSecure := secure.New(secure.Config(cfg.Secure))
		go appSecure.MustGetPermissionsNumbers(permissionsChan)
	}

	metrics := prometheusMetrics.MustCreate(&cfg.Prometheus)
	domainService := service.New(mysql.WithRepository(&cfg.Storage),
		service.WithMetrics(metrics))

	if cfg.UseKafka {
		kafka.MustRun(domainService, &cfg.Kafka, cfg.Instance)
	}

	server := restServer.MustCreate(&cfg.HttpServer, cfg.QueryTimeout, domainService, metrics, cfg.Env,
		cfg.Signature, permissionsChan)
	server.MustRun()

	var viewer *db_reader.Reader
	if cfg.Storage.ViewerPort != 0 && (cfg.Env == config.EnvironmentLocal || cfg.Env == config.EnvironmentDebug) {
		viewer = db_reader.New(domainService.SQLRepository.DB(), cfg.Storage.ViewerPort)
		viewer.Start()
	}

	defer func(repo repository.SQLDBInterface) {
		if repo != nil {
			_ = repo.Close()
		}
	}(domainService.SQLRepository)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	fmt.Println() // Так красивее, если вывод логов производится в стандартный терминал
	slog.Info(fmt.Sprintf("%s signal received. Shutdown started", sig))

	server.Shutdown()

	if viewer != nil {
		viewer.Shutdown()
	}
}

func clearScreen() error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
