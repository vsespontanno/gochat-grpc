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
		"bootstrap.servers": broker,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{consumer: consumer, grpcClient: grpcClient}, nil
}

func (p *KafkaProducer) Produce(topic, value string) error {
	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(value),
	}, nil)
	if err != nil {
		log.Println("Error while producing msg: ", err)
		return err
	}

	p.producer.Flush(15 * 1000)
	return nil
}

func (c *KafkaConsumer) Subscribe(topic string) error {
	err := c.consumer.Subscribe(topic, nil)
	if err != nil {
		return err
	}
	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received message: %s\n", string(msg.Value))
		} else {
			fmt.Printf("Error while consuming msg: %v\n", err)
		}

	}
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
