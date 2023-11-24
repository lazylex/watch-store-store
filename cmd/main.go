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
	"os"
	"os/signal"
)

func main() {
	cfg := config.MustLoad()
	log := logger.MustCreate(cfg.Env, cfg.Instance)
	metrics := prometheusMetrics.MustCreate(&cfg.Prometheus, log)
	domainService := service.New(mysql.WithRepository(&cfg.Storage, log), service.WithLogger(log),
		service.WithMetrics(metrics))

	server := restServer.New(&cfg.HttpServer, cfg.QueryTimeout, domainService, log, metrics, cfg.Env, cfg.Signature)
	server.MustRun()

	if cfg.UseKafka {
		kafka.MustRun(domainService, &cfg.Kafka, log, cfg.Instance)
	}

	defer func(Repository repository.Interface) {
		_ = Repository.Close()
	}(domainService.Repository)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	fmt.Println() // так красивее, если вывод логов производится в стандартный терминал
	log.Info(fmt.Sprintf("%s signal received. Shutdown started", sig))

	server.Shutdown()
}
