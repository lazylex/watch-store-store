package kafka

import (
	"fmt"
	"github.com/lazylex/watch-store-store/internal/adapters/message_broker/kafka/consumer/request_count"
	"github.com/lazylex/watch-store-store/internal/adapters/message_broker/kafka/consumer/update_price"
	"github.com/lazylex/watch-store-store/internal/adapters/message_broker/kafka/producer/response_count"
	"github.com/lazylex/watch-store-store/internal/config"
	"github.com/lazylex/watch-store-store/internal/dto"
	internalLogger "github.com/lazylex/watch-store-store/internal/logger"
	"github.com/lazylex/watch-store-store/internal/ports/service"
	"log/slog"
	"os"
)

const countChannelBufferSize = 10

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

	if len(cfg.RequestCountTopic) > 0 && len(cfg.ResponseCountTopic) > 0 {
		countCh := make(chan dto.ArticleAmount, countChannelBufferSize)
		go request_count.ListenTopic(service, cfg.Brokers, cfg.RequestCountTopic, instance, countCh)
		go response_count.Serve(cfg.Brokers, cfg.ResponseCountTopic, instance, countCh)

		topicsInService += 2
	} else {
		log.Error("not configured Kafka count topics")
	}

	if topicsInService > 0 {
		log.Info(fmt.Sprintf("kafka topics in service: %d", topicsInService))
	} else {
		log.Info("kafka: no topics to service")
	}

}
