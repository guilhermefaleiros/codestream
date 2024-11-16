package kafka

import (
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type MessageProcessor interface {
	Process(msg *ckafka.Message) error
}

type Consumer struct {
	configMap      *ckafka.ConfigMap
	topic          string
	messageChannel chan *ckafka.Message
	processor      MessageProcessor
	done           chan bool
}

func NewConsumer(configMap *ckafka.ConfigMap, topic string, processor MessageProcessor) *Consumer {
	return &Consumer{
		configMap:      configMap,
		topic:          topic,
		messageChannel: make(chan *ckafka.Message, 100),
		processor:      processor,
		done:           make(chan bool),
	}
}

func (c *Consumer) Start() {
	go func() {
		consumer, err := ckafka.NewConsumer(c.configMap)
		if err != nil {
			log.Fatalf("Failed to create Kafka consumer: %v", err)
		}
		defer func(consumer *ckafka.Consumer) {
			err := consumer.Close()
			if err != nil {
				log.Fatalf("Failed to close Kafka consumer: %v", err)
			}
		}(consumer)

		err = consumer.SubscribeTopics([]string{c.topic}, nil)
		if err != nil {
			log.Fatalf("Failed to subscribe to topic %s: %v", c.topic, err)
		}

		log.Printf("Consumer started for topic: %s", c.topic)

		for {
			select {
			case <-c.done:
				log.Println("Stopping consumer...")
				return
			default:
				msg, err := consumer.ReadMessage(-1)
				if err == nil {
					c.messageChannel <- msg
				} else {
					log.Printf("Error reading message: %v", err)
				}
			}
		}
	}()
}

func (c *Consumer) ProcessMessages() {
	go func() {
		for {
			select {
			case <-c.done:
				log.Println("Stopping message processing...")
				return
			case msg := <-c.messageChannel:
				if err := c.processor.Process(msg); err != nil {
					log.Printf("Error processing message: %v", err)
				}
			}
		}
	}()
}

func (c *Consumer) Stop() {
	close(c.done)
	close(c.messageChannel)
}
