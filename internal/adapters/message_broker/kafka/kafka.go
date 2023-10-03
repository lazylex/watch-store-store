package kafka

import (
	"github.com/lazylex/watch-store/store/internal/adapters/message_broker/kafka/consumer/update_price"
	"github.com/lazylex/watch-store/store/internal/config"
	"github.com/lazylex/watch-store/store/internal/ports/service"
	"log/slog"
)

func Start(service service.Interface, cfg *config.Config, logger *slog.Logger) {
	go update_price.UpdatePrice(service, logger, cfg.Brokers)
}
