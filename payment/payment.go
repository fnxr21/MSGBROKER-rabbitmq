package main

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	common "github/com/fnxr21/msgbroker-common"
	"log"
)

var (
	amqpUser = "guest"
	amqpPass = "guest"
	amqpHost = "192.168.100.132" 
	amqpPort = "5672"
)

func main() {

	ch, close := common.ConnectAmqp(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	listen(ch)
}

func listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(common.OrderCreatedEvent, true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct {
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			o := &common.Order{}
			if err := json.Unmarshal(d.Body, o); err != nil {
				d.Nack(false, false)
				log.Printf("failed to unmarshal order:%v", err)
				continue
			}
			paymentLink, err := createPaymentLink()
			if err != nil {
				log.Printf("failed to create payment: %v", err)
				continue
			}
			log.Printf("Payment link generated: %s", paymentLink)
		}
	}()

	log.Printf("AMQP Listening. To exist press  CTRL+C")

	<-forever
}

func createPaymentLink() (string, error) {
	return "dummy-payment-link.com", nil
}
