package main

import (
	"fmt"
	"log"

	"github.com/spandan12/rabitmqreceive/pkg/rabbitmq"
)

const (
	queue      string = "send_cleaner_worker"
	bindingKey string = "signoi.analysis.dataset.csv.chunk.cleaned"
)

func main() {
	url := "amqp://guest:guest@localhost:5672/"
	if err := rabbitMQ(url); err != nil {
		fmt.Println("error:", err)
	}
}

// RabbitMQ
func rabbitMQ(url string) error {
	rbm, err := rabbitmq.New(&rabbitmq.Config{
		URL:          url,
		Exchange:     "analysis",
		ExchangeType: "topic",
	})
	if err != nil {
		return err
	}
	defer rbm.Close()

	if err := rbm.NewExchange(); err != nil {
		return err
	}

	if err := rbm.NewQueue(queue, bindingKey); err != nil {
		return err
	}

	msgs, err := rbm.Consume(queue)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
	return nil
}
