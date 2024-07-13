package request_count

import (
	"context"
	"errors"
	"fmt"
	"github.com/lazylex/watch-store-store/internal/domain/value_objects/article"
	"github.com/lazylex/watch-store-store/internal/dto"
	"github.com/lazylex/watch-store-store/internal/logger"
	"github.com/lazylex/watch-store-store/internal/ports/repository"
	"github.com/lazylex/watch-store-store/internal/ports/service"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"time"
)

const attemptsUntilAlarm = 6

// ListenTopic ожидает сообщение в топике с запросом количества товара (строка с артикулом), получает это количество из
// сервисного слоя и отправляет в канал countChan.
func ListenTopic(service service.Interface, brokers []string, topic, instance string, countChan chan<- uint) {
	var err error
	var m kafka.Message
	var attempts int
	var amount uint

	ctx := context.Background()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		MaxBytes: 10e6,
		GroupID:  instance,
	})

	log := slog.With(slog.String(logger.OPLabel, "kafka.consumer.request_count.ListenTopic"))
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

		data := dto.ArticleDTO{Article: article.Article(m.Value)}
		err = data.Validate()

		if err != nil {
			log.Warn(err.Error())
		} else {
			log.Info(fmt.Sprintf("reading product count (article %s)", data.Article))
			if amount, err = service.AmountInStock(ctx, data); err != nil {
				if attempts < attemptsUntilAlarm {
					log.Warn(err.Error())
				} else {
					log.Error(err.Error())
				}

				if errors.Is(err, repository.ErrNoRecord) {
					canFetchMessage = true
					amount = 0
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
				countChan <- amount
			}
		}

	}

	if err = r.Close(); err != nil {
		log.Error("failed to close reader:", err)
	}
}
