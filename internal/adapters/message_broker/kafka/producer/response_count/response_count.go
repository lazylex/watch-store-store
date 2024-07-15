package response_count

import (
	"context"
	"fmt"
	"github.com/lazylex/watch-store-store/internal/dto"
	"github.com/lazylex/watch-store-store/internal/logger"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

const attemptsUntilAlarm = 6

// Serve прослушивает канал countChan и отправляет его содержимое в формате JSON в топик topic. При этом к данным из
// канала добавляется поле с названием экземпляра приложения instance.
func Serve(brokers []string, topic, instance string, countChan <-chan dto.ArticleAmount) {
	log := slog.With(slog.String(logger.OPLabel, "kafka.consumer.response_count.Serve"))
	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		MaxAttempts:            attemptsUntilAlarm,
		AllowAutoTopicCreation: true,
	}

	for {
		c, ok := <-countChan

		if !ok {
			if err := w.Close(); err != nil {
				log.Error("failed to close writer:", err)
			}
		}

		if err := w.WriteMessages(context.Background(),
			kafka.Message{
				Value: []byte(fmt.Sprintf("{\"instance\":%s,\"article\":%s,\"count\":%d}", instance, c.Article, c.Amount)),
			},
		); err != nil {
			log.Error("failed to write messages:" + err.Error())
		} else {
			log.Info(fmt.Sprintf("quantity of goods with article %s was successfully sent", c.Article))
		}
	}
}
