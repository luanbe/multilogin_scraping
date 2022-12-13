package helper

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type RabbitMQBroker interface {
	PublishMessage(exchangeType, exchangeName, routingKey, body string)
	ConsumeMessage(exchangeType, exchangeName, queueName, routingKey string)
}
type RabbitMQHelper struct {
	ServerURL string
	Connect   *amqp.Connection
	Channel   *amqp.Channel
}

func NewRabbitMQ(serverURL string) RabbitMQBroker {
	r := &RabbitMQHelper{ServerURL: serverURL}
	r.CreateRabbitMQ()
	return r
}

func (r *RabbitMQHelper) CreateRabbitMQ() {
	conn, err := amqp.Dial(r.ServerURL)
	failOnError(err, "Failed to connect to RabbitMQHelper")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	r.Connect = conn
	r.Channel = ch

}

func (r *RabbitMQHelper) PublishMessage(exchangeType, exchangeName, routingKey, body string) {
	defer r.Connect.Close()
	defer r.Channel.Close()

	err := r.Channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = r.Channel.PublishWithContext(ctx,
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

}
func (r *RabbitMQHelper) ConsumeMessage(exchangeType, exchangeName, queueName, routingKey string) {
	defer r.Connect.Close()
	defer r.Channel.Close()
	err := r.Channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := r.Channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		true,      // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = r.Channel.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	messages, err := r.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// Build a welcome message.
	log.Println("Successfully connected to RabbitMQHelper")
	log.Println("Waiting for messages")

	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)

	go func() {
		for message := range messages {
			// For example, show received message in a console.
			log.Printf(" > Received message: %s\n", message.Body)
		}
	}()

	<-forever

}
