package kafka

import (
	"github.com/lazylex/watch-store/store/internal/adapters/message_broker/kafka/consumer/update_price"
	"github.com/lazylex/watch-store/store/internal/config"
	"github.com/lazylex/watch-store/store/internal/ports/service"
	"log/slog"
)

// Start предназначен для запуска консьюмеров/продюсеров Кафки. Если в конфигурации cfg не задано имя топика, то
// соответствующий ему консьюмер/продюсер не будет запущен. Работа приложения будет продолжена
func Start(service service.Interface, cfg *config.Config, logger *slog.Logger) {
	if len(cfg.UpdatePriceTopic) > 0 {
		go update_price.UpdatePrice(service, logger, cfg)
	} else {
		logger.Error("not configured Kafka Update Price topic")
	}
}
