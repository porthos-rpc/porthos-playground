package collector

import (
	"fmt"

	"github.com/porthos-rpc/porthos-playground/models"
	"github.com/porthos-rpc/porthos-playground/storage"
	"github.com/streadway/amqp"
)

const (
	specsQueueName = "porthos.specs"
)

// Collector is responsible for consuming specs from the broker and persist them through storage.
type Collector struct {
	broker          *amqp.Connection
	channel         *amqp.Channel
	deliveryChannel <-chan amqp.Delivery
	storage         storage.Storage
}

// NewCollector creates a new specs collector.
func NewCollector(brokerURL string, s storage.Storage) *Collector {
	broker, err := NewBroker(brokerURL)

	if err != nil {
		panic(err)
	}

	ch, err := broker.Channel()

	if err != nil {
		panic(err)
	}

	_, err = ch.QueueDeclare(
		specsQueueName, // name
		true,           // durable
		false,          // delete when usused
		false,          // exclusive
		false,          // noWait
		nil,            // arguments
	)

	dc, _ := ch.Consume(
		specsQueueName, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)

	return &Collector{broker, ch, dc, s}
}

// Start the specs collector.
func (c *Collector) Start() {
	for d := range c.deliveryChannel {
		specs, err := models.UnmarshalSpecs(d.Body)

		if err != nil {
			fmt.Errorf("Error parsing specs %s.", err)
			return
		}

		c.storage.SaveServiceSpecs(specs)
	}
}

// Stop the specs collector and release resources.
func (c *Collector) Stop() {
	c.channel.Close()
	c.broker.Close()
}
