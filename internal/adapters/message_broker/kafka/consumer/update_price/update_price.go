package update_price

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lazylex/watch-store/store/internal/dto"
	"github.com/lazylex/watch-store/store/internal/logger"
	"github.com/lazylex/watch-store/store/internal/ports/service"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

// UpdatePrice обновляет цену товара, находящегося в продаже, если считывает в топике store.update-price новую цену
func UpdatePrice(service service.Interface, log *slog.Logger, brokers []string) {
	r := kafka.NewReader(kafka.ReaderConfig{Brokers: brokers, Topic: "store.update-price", Partition: 0, MaxBytes: 10e6})

	log = log.With(slog.String(logger.OPLabel, "kafka.consumer.UpdatePrice"))
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}

		var transferObject dto.ArticleWithPriceDTO
		err = json.Unmarshal(m.Value, &transferObject)

		if err != nil {
			log.Warn("error unmarshal JSON")
		} else {
			err = transferObject.Validate()
			if err != nil {
				log.Warn(err.Error())
			} else {
				log.Info(fmt.Sprintf("reading updating price to %.2f (article %s)",
					transferObject.Price, transferObject.Article))
				if err = service.ChangePriceInStock(context.Background(), transferObject); err != nil {
					log.Warn(err.Error())
				}
			}
		}
	}

	if err := r.Close(); err != nil {
		log.Error("failed to close reader:", err)
	}
}
