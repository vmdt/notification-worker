package rabbitmq

import (
	"context"

	"github.com/streadway/amqp"
	"github.com/vmdt/notification-worker/pkg/logger"
)

type IConsumer interface {
	ConsumeMessage(msg interface{}, name string, queue string, key string) error
}

type Consumer struct {
	cfg     *RabbitMQConfig
	conn    *amqp.Connection
	log     logger.ILogger
	handler func(queue string, msg amqp.Delivery) error
	ctx     context.Context
}

func (c Consumer) ConsumeMessage(msg interface{}, name string, queue string, key string) error {
	channel, err := c.conn.Channel()
	if err != nil {
		c.log.Error("Error in opening channel to consume message")
	}
	defer channel.Close()

	err = channel.ExchangeDeclare(
		name,       // name
		c.cfg.Kind, // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		c.log.Error("Error in declaring exchange")
	}

	q, err := channel.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		c.log.Error("Error in declaring queue")
		return err
	}

	err = channel.QueueBind(
		q.Name, // queue name
		key,    // routing key
		name,   // exchange
		false,  // no-wait
		nil,    // arguments
	)

	if err != nil {
		c.log.Error("Error in binding queue")
		return err
	}

	deliveries, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		c.log.Error("Error in consuming message")
		return err
	}

	go func() {
		for {
			select {
			case <-c.ctx.Done():
				defer func(ch *amqp.Channel) {
					err := ch.Close()
					if err != nil {
						c.log.Error("Error in closing channel")
					}
				}(channel)
				c.log.Info("RabbitMQ channel is closed")
				return
			case delivery, ok := <-deliveries:
				{
					if !ok {
						c.log.Error("Error in consuming message")
						return
					}
					c.log.Infof("Received message: %s", delivery.Body)
					err := c.handler(queue, delivery)
					if err != nil {
						c.log.Error(err.Error())
					}

					err = delivery.Ack(false)
					if err != nil {
						c.log.Error("Error in acknowledging message")
					}
				}
			}
		}
	}()
	c.log.Infof("Waiting for messages in queue: %s", q.Name)
	return nil
}

func NewConsumer(ctx context.Context, cfg *RabbitMQConfig, conn *amqp.Connection, log logger.ILogger, handler func(queue string, msg amqp.Delivery) error) IConsumer {
	return &Consumer{ctx: ctx, cfg: cfg, conn: conn, log: log, handler: handler}
}
