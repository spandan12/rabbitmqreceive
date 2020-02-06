package rabbitmq

import (
	"github.com/streadway/amqp"
)

// Config ...
type Config struct {
	URL          string
	Exchange     string
	ExchangeType string
}

// Rabbitmq ...
type Rabbitmq struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchange     string
	exchangetype string
}

// New ...
func New(config *Config) (*Rabbitmq, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Rabbitmq{
		conn:         conn,
		channel:      ch,
		exchange:     config.Exchange,
		exchangetype: config.ExchangeType,
	}, nil
}

// NewExchange ...
func (r *Rabbitmq) NewExchange() error {
	return r.channel.ExchangeDeclare(
		r.exchange,     // name
		r.exchangetype, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
}

// NewQueue ...
func (r *Rabbitmq) NewQueue(queue string, bindingKey string) error {
	q, err := r.channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = r.channel.QueueBind(
		q.Name,     // queue name
		bindingKey, // routing key
		r.exchange, // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

// Publish ...
func (r *Rabbitmq) Publish(routingKey string, contentType string, body []byte) error {
	return r.channel.Publish(
		r.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		})
}

// Consume ...
func (r *Rabbitmq) Consume(queue string) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queue, // queue
		"",    // consumer
		true,  // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)
}

// Close ...
func (r *Rabbitmq) Close() {
	r.conn.Close()
	r.channel.Close()
}
