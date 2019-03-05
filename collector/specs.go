package collector

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/porthos-rpc/porthos-playground/models"
	"github.com/porthos-rpc/porthos-playground/storage"
	"github.com/sirupsen/logrus"
	amqplib "github.com/streadway/amqp"
)

const specsQueueName = "porthos.specs"

var consumerTagSeq uint64

// Collector is responsible for consuming specs from the broker and persist them through storage.
type Collector struct {
	conn    *connection
	storage storage.Storage
	m       sync.RWMutex
	closed  bool
}

// NewCollector creates a new specs collector.
func NewCollector(brokerURL string, s storage.Storage) *Collector {
	conn, err := NewConnection(brokerURL)

	if err != nil {
		logrus.WithError(err).Fatal("Failed to create a new connection.")
	}

	return &Collector{
		conn:    conn,
		storage: s,
	}
}

func createUniqueConsumerTagName() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return fmt.Sprintf("ctag-%s-%s-%d", hostname, os.Args[0], atomic.AddUint64(&consumerTagSeq, 1))
}

func (c *Collector) doConsume() error {
	logrus.WithField("queue", specsQueueName).Debug("Setting up consumer channel...")

	channel, err := c.conn.openChannel()

	if err != nil {
		return err
	}

	defer channel.Close()

	err = c.setupTopology(channel)

	if err != nil {
		return err
	}

	msgs, err := channel.Consume(
		specsQueueName,                // queue
		createUniqueConsumerTagName(), // consumer
		true,                          // auto ack
		false,                         // exclusive
		false,                         // no local
		false,                         // no wait
		nil,                           // args
	)

	if err != nil {
		return err
	}

	logrus.WithField("queue", specsQueueName).Info("Consuming messages...")

	for d := range msgs {
		specs, err := models.UnmarshalSpecs(d.Body)

		if err != nil {
			logrus.WithError(err).Error("Failed to parse specs.")
			return err
		}

		c.storage.SaveServiceSpecs(specs)
	}

	return nil
}

func (c *Collector) setupTopology(channel *amqplib.Channel) error {
	_, err := channel.QueueDeclare(
		specsQueueName, // name
		true,           // durable
		false,          // delete when usused
		false,          // exclusive
		false,          // noWait
		nil,            // arguments
	)

	return err
}

// Start the specs collector.
func (c *Collector) Start() {
	rs := c.conn.NotifyReestablish()

	for !c.closed {
		if !c.conn.IsConnected() {
			logrus.Info("Connection not established. Waiting connection to be reestablished.")

			<-rs

			continue
		}

		err := c.doConsume()

		if err == nil {
			logrus.WithFields(logrus.Fields{
				"queue":  specsQueueName,
				"closed": c.closed,
			}).Info("Consumption finished.")
		} else {
			logrus.WithFields(logrus.Fields{
				"queue": specsQueueName,
				"error": err,
			}).Error("Error consuming events.")
		}
	}
}

// Close closes the specs collector and release resources.
func (c *Collector) Close() {
	func() {
		c.m.Lock()
		defer c.m.Unlock()
		c.closed = true
	}()

	c.conn.Close()
}
