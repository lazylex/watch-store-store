package kafka

import (
	"fmt"
	"github.com/lazylex/watch-store/store/internal/adapters/message_broker/kafka/consumer/update_price"
	"github.com/lazylex/watch-store/store/internal/config"
	internalLogger "github.com/lazylex/watch-store/store/internal/logger"
	"github.com/lazylex/watch-store/store/internal/ports/service"
	"log/slog"
	"os"
)

// MustRun предназначен для запуска consumers/producers Кафки. Если в конфигурации cfg не задано имя топика, то
// соответствующий ему consumer/producer не будет запущен. Работа приложения будет продолжена.
func MustRun(service service.Interface, cfg *config.Kafka, instance string) {
	var topicsInService int
	log := slog.With(slog.String(internalLogger.OPLabel, "kafka.MustRun"))

	if len(cfg.Brokers) < 1 {
		log.Error("empty kafka brokers list")
		os.Exit(1)
	}

	if len(cfg.UpdatePriceTopic) > 0 {
		go update_price.UpdatePrice(service, cfg.Brokers, cfg.UpdatePriceTopic, instance)
		topicsInService++
	} else {
		log.Error("not configured Kafka Update Price topic")
	}

	if topicsInService > 0 {
		log.Info(fmt.Sprintf("kafka topics in service: %d", topicsInService))
	} else {
		log.Info("kafka: no topics to service")
	}

}
