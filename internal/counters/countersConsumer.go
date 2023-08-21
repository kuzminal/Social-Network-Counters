package counters

import (
	"SocialNetCounters/internal/helper"
	"SocialNetCounters/internal/store"
	"SocialNetCounters/models"
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"os"
)

type Consumer struct {
	Reader        *kafka.Reader
	CountersStore store.CountersStore
}

func NewCountersConsumer(countersStore store.CountersStore) Consumer {
	brokerHost := helper.GetEnvValue("KAFKA_BROKER_HOST", "localhost")
	brokerPort := helper.GetEnvValue("KAFKA_BROKER_PORT", "9092")
	l := log.New(os.Stdout, "kafka reader: ", 0)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerHost + ":" + brokerPort},
		Topic:   "counters",
		GroupID: "counters",
		Logger:  l,
	})
	return Consumer{Reader: r, CountersStore: countersStore}
}

func (c *Consumer) ReadCountersInfo(ctx context.Context) {
	for {
		msg, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		var messageInfo models.Message
		err = msgpack.Unmarshal(msg.Value, &messageInfo)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("received: ", messageInfo.FromUser)
		if messageInfo.IsRead {
			_, err = c.CountersStore.DecrCounter(messageInfo.ToUser)
			if err != nil {
				log.Println(err)
			}
		} else {
			_, err = c.CountersStore.IncrCounter(messageInfo.ToUser)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
