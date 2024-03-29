package main

import (
	"fmt"
	"github.com/lazylex/watch-store/store/internal/adapters/message_broker/kafka"
	restServer "github.com/lazylex/watch-store/store/internal/adapters/rest/server"
	"github.com/lazylex/watch-store/store/internal/config"
	"github.com/lazylex/watch-store/store/internal/logger"
	prometheusMetrics "github.com/lazylex/watch-store/store/internal/metrics"
	"github.com/lazylex/watch-store/store/internal/ports/repository"
	"github.com/lazylex/watch-store/store/internal/repository/mysql"
	"github.com/lazylex/watch-store/store/internal/service"
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
	metrics := prometheusMetrics.MustCreate(&cfg.Prometheus)
	domainService := service.New(mysql.WithRepository(&cfg.Storage),
		service.WithMetrics(metrics))

	server := restServer.MustCreate(&cfg.HttpServer, cfg.QueryTimeout, domainService, metrics, cfg.Env,
		cfg.Signature)
	server.MustRun()

	if cfg.UseKafka {
		kafka.MustRun(domainService, &cfg.Kafka, cfg.Instance)
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
	fmt.Println() // так красивее, если вывод логов производится в стандартный терминал
	slog.Info(fmt.Sprintf("%s signal received. Shutdown started", sig))

	server.Shutdown()
}

func clearScreen() error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
