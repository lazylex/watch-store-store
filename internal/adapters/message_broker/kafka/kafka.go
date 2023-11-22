package kafka

import (
	"fmt"
	"github.com/lazylex/watch-store/store/internal/adapters/message_broker/kafka/consumer/update_price"
	"github.com/lazylex/watch-store/store/internal/config"
	"github.com/lazylex/watch-store/store/internal/logger"
	"github.com/lazylex/watch-store/store/internal/ports/service"
	"log/slog"
	"os"
)

// MustRun предназначен для запуска консьюмеров/продюсеров Кафки. Если в конфигурации cfg не задано имя топика, то
// соответствующий ему консьюмер/продюсер не будет запущен. Работа приложения будет продолжена
func MustRun(service service.Interface, cfg *config.Kafka, log *slog.Logger, instance string) {
	var topicsInService int
	locLog := log.With(slog.String(logger.OPLabel, "kafka.MustRun"))

	if len(cfg.Brokers) < 1 {
		locLog.Error("empty kafka brokers list")
		os.Exit(1)
	}

	if len(cfg.UpdatePriceTopic) > 0 {
		go update_price.UpdatePrice(service, log, cfg.Brokers, cfg.UpdatePriceTopic, instance)
		topicsInService++
	} else {
		locLog.Error("not configured Kafka Update Price topic")
	}

	if topicsInService > 0 {
		locLog.Info(fmt.Sprintf("kafka topics in service: %d", topicsInService))
	} else {
		locLog.Info("kafka: no topics to service")
	}

}
