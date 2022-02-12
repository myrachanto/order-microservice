package events

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var Producer *kafka.Producer

func SetupProducer() {
	var err error
	Producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost",
		"security.prototcol": "SASL_SSL",
		"sasl.username":      "",
		"sasl.password":      "",
		"sasl.mechanism":     "PLAIN",
	})

	if err != nil {
		log.Println("producer error")
	}
	defer Producer.Close()
}

func Produce(topic, key string, message interface{}) {
	value, _ := json.Marshal(message)
	Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
	}, nil)
	Producer.Flush(15000)
}
