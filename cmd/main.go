package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	restHandles "github.com/lazylex/watch-store/store/internal/adapters/rest/handlers"
	"github.com/lazylex/watch-store/store/internal/adapters/rest/router"
	"github.com/lazylex/watch-store/store/internal/config"
	"github.com/lazylex/watch-store/store/internal/logger"
	"github.com/lazylex/watch-store/store/internal/ports/repository"
	"github.com/lazylex/watch-store/store/internal/repository/mysql"
	"github.com/lazylex/watch-store/store/internal/service"
	"net/http"
	"os"
	"os/signal"
)

var configPath = flag.String("config", "", "путь к файлу конфигурации")

func main() {
	flag.Parse()
	cfg := config.MustLoad(configPath)
	log := logger.MustCreate(cfg.Env, cfg.Instance)
	domainService := service.New(
		mysql.WithRepository(cfg, log),
		service.WithLogger(log),
	)
	handlers := restHandles.New(domainService, log, cfg.QueryTimeout)

	srv := &http.Server{
		Handler:      router.New(cfg, log, handlers),
		Addr:         cfg.Address,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		log.Info("start http server on " + cfg.Address)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server startup error")
			os.Exit(1)
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

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to gracefully shutdown http server")
	} else {
		log.Info("gracefully shut down http server")
	}
}
