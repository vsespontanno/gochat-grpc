package messaging

import (
	"fmt"
	"log"

	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/vsespontanno/gochat-grpc/internal/client"
	"github.com/vsespontanno/gochat-grpc/internal/models"
)

type KafkaProducer struct {
	producer *kafka.Producer
}

type KafkaConsumer struct {
	consumer   *kafka.Consumer
	grpcClient client.GRPCClient
}

func NewKafkaProducer(broker string) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer: producer}, nil
}

func NewKafkaConsumer(broker, groupID string, grpcClient client.GRPCClient) (*KafkaConsumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  broker,
		"group.id":           groupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{consumer: consumer, grpcClient: grpcClient}, nil
}

func (p *KafkaProducer) Produce(topic string, msg models.Message) error {
	value, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return err
	}
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, nil)

	if err != nil {
		log.Println("Error while producing msg: ", err)
		return err
	}

	p.producer.Flush(15 * 1000)
	return nil
}

func (c *KafkaConsumer) Subscribe(topic string) {
	err := c.consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Failed to consume: %v", err)
	}
	go func() {
		for {
			msg, err := c.consumer.ReadMessage(-1)
			if err == nil {
				var message models.Message
				err := json.Unmarshal(msg.Value, &message)
				if err != nil {
					log.Fatalf("Failed to unmarshal message: %v", err)
				}
				fmt.Printf("Received message: %s\n", message)
			} else {
				log.Fatalf("Error while consuming msg: %v\n", err)
			}
		}
	}()
}

func (c *KafkaConsumer) ProcessMessage(msg []byte) (bool, error) {
	var message models.Message
	err := json.Unmarshal(msg, &message)
	if err != nil {
		return false, err
	}

	if message.Recipient != "me" {
		return false, nil
	}

	///....

	return true, nil
}
