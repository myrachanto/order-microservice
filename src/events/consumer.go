package events

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func init() {
	log.SetPrefix("Kafka :")
}
func main() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost",
		"security.prototcol": "SASL_SSL",
		"sasl.username":      "",
		"sasl.password":      "",
		"sasl.mechanism":     "PLAIN",
		"group.id":           "mygroup",
		"auto.offset.reset":  "earliest",
	})
	if err != nil {
		log.Println("Error : something went wrong Initilalizing kafka")
	}
	consumer.SubscribeTopics([]string{"user_topic"}, nil)
	for err != nil {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
			return
		}
		fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		var message map[string]interface{}
		json.Unmarshal(msg.Value, &message)
	}
	consumer.Close()
}
