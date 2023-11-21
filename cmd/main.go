package main

import (
	"flag"
	"fmt"
	"github.com/lazylex/watch-store/store/internal/adapters/message_broker/kafka"
	restServer "github.com/lazylex/watch-store/store/internal/adapters/rest/server"
	"github.com/lazylex/watch-store/store/internal/config"
	"github.com/lazylex/watch-store/store/internal/logger"
	prometheusMetrics "github.com/lazylex/watch-store/store/internal/metrics"
	"github.com/lazylex/watch-store/store/internal/ports/repository"
	"github.com/lazylex/watch-store/store/internal/repository/mysql"
	"github.com/lazylex/watch-store/store/internal/service"
	"os"
	"os/signal"
)

var configPath = flag.String("config", "", "путь к файлу конфигурации")

func main() {
	flag.Parse()
	cfg := config.MustLoad(configPath)
	log := logger.MustCreate(cfg.Env, cfg.Instance)
	metrics := prometheusMetrics.MustCreate(cfg, log)

	domainService := service.New(
		mysql.WithRepository(cfg, log),
		service.WithLogger(log),
		service.WithMetrics(metrics),
	)

	server := restServer.New(
		cfg.Address, cfg.ReadTimeout, cfg.WriteTimeout, cfg.IdleTimeout, cfg.ShutdownTimeout, cfg.QueryTimeout,
		domainService, log, metrics, cfg.Env, cfg.Signature)

	server.MustRun()

	go func() {
		if cfg.UseKafka {
			if len(cfg.Brokers) > 0 {
				kafka.Start(domainService, cfg, log)
			} else {
				log.Error("empty kafka brokers list")
			}
		}
	}()

	defer func(Repository repository.Interface) {
		err := Repository.Close()
		if err != nil {
			log.Error("error close repository")
		}
		log.Info("close repository")
	}(domainService.Repository)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	fmt.Println() // так красивее, если вывод логов производится в стандартный терминал
	log.Info(fmt.Sprintf("%s signal received. Shutdown started", sig))

	server.Shutdown()
}
