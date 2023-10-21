package update_price

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lazylex/watch-store/store/internal/config"
	"github.com/lazylex/watch-store/store/internal/dto"
	"github.com/lazylex/watch-store/store/internal/logger"
	"github.com/lazylex/watch-store/store/internal/ports/repository"
	"github.com/lazylex/watch-store/store/internal/ports/service"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"time"
)

const attemptsUntilAlarm = 6

// UpdatePrice обновляет цену товара, находящегося в продаже, если считывает в топике store.update-price новую цену.
// Автокоммит не выполняется. При ошибке обновления цены смещение в Кафке не сохраняется, а производятся новые попытки
// обновления. Каждая последующая попытка производится через период, на десять секунд дольше предыдущего. Через
// attemptsUntilAlarm попыток, в лог выводится ошибка, а не предупреждение
func UpdatePrice(service service.Interface, log *slog.Logger, cfg *config.Config) {
	var err error
	var m kafka.Message
	var attempts int
	ctx := context.Background()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.UpdatePriceTopic,
		MaxBytes: 10e6,
		GroupID:  cfg.Instance,
	})

	log = log.With(slog.String(logger.OPLabel, "kafka.consumer.UpdatePrice"))
	canFetchMessage := true
	for {
		if canFetchMessage {
			m, err = r.FetchMessage(ctx)
			if err != nil {
				break
			}
			attempts = 0
		}
		canFetchMessage = true

		var data dto.ArticleWithPriceDTO
		err = json.Unmarshal(m.Value, &data)

		if err != nil {
			log.Warn("error unmarshal JSON")
		} else {
			err = data.Validate()
			if err != nil {
				log.Warn(err.Error())
			} else {
				log.Info(fmt.Sprintf("reading updating price to %.2f (article %s)", data.Price, data.Article))
				if err = service.ChangePriceInStock(ctx, data); err != nil {
					if attempts < attemptsUntilAlarm {
						log.Warn(err.Error())
					} else {
						log.Error(err.Error())
					}

					if errors.Is(err, repository.ErrNoRecord) {
						canFetchMessage = true
					} else {
						canFetchMessage = false
						attempts++
						time.Sleep(time.Second * time.Duration(10*attempts))
					}
				}

				if canFetchMessage {
					err = r.CommitMessages(ctx, m)
					if err != nil {
						log.Warn(err.Error())
					}
				}
			}
		}
	}

	if err = r.Close(); err != nil {
		log.Error("failed to close reader:", err)
	}
}
