package helper

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"
)

type RabbitMQBroker interface {
	PublishMessage(exchangeType, exchangeName, routingKey string, body []byte) error
	ConsumeMessage(exchangeType, exchangeName, queueName, routingKey string) (<-chan amqp.Delivery, *RabbitMQHelper)
	CloseRabbitMQ()
}
type RabbitMQHelper struct {
	ServerURL string
	Connect   *amqp.Connection
	Channel   *amqp.Channel
	Logger    *zap.Logger
	Utils     *UtilHelper
}

func NewRabbitMQ(serverURL string, rabbitLog *zap.Logger) RabbitMQBroker {
	utils := &UtilHelper{}
	r := &RabbitMQHelper{ServerURL: serverURL, Logger: rabbitLog, Utils: utils}
	r.CreateRabbitMQ()
	return r
}

func (r *RabbitMQHelper) CreateRabbitMQ() {
	conn, err := amqp.Dial(r.ServerURL)
	r.Utils.failOnError(err, "Failed to connect to RabbitMQHelper")

	ch, err := conn.Channel()
	r.Utils.failOnError(err, "Failed to open a channel")

	r.Connect = conn
	r.Channel = ch

}

func (r *RabbitMQHelper) CloseRabbitMQ() {
	r.Connect.Close()
	r.Channel.Close()
}
func (r *RabbitMQHelper) PublishMessage(exchangeType, exchangeName, routingKey string, body []byte) error {
	if err := r.Channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.Channel.PublishWithContext(ctx,
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		}); err != nil {
		return err
	}
	return nil
}
func (r *RabbitMQHelper) ConsumeMessage(exchangeType, exchangeName, queueName, routingKey string) (<-chan amqp.Delivery, *RabbitMQHelper) {
	err := r.Channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	r.Utils.failOnError(err, "Failed to declare an exchange")

	q, err := r.Channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		true,      // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	r.Utils.failOnError(err, "Failed to declare a queue")

	err = r.Channel.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil)
	r.Utils.failOnError(err, "Failed to bind a queue")

	messages, err := r.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	r.Utils.failOnError(err, "Failed to register a consumer")

	// Build a welcome message.
	r.Logger.Info("Successfully connected to RabbitMQHelper")
	r.Logger.Info("Waiting for messages")

	return messages, r
}
