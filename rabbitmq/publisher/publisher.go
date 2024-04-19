package publisher

import (
	"context"
	"errors"
	"log"
	"message_queue_system/domain"

	"github.com/streadway/amqp"
)

type PublishTaskRequest struct {
	ReqBytes     []byte `json:"reqBytes"`
	QueueName    string `json:"queueName"`
	ExchangeName string `json:"exchangeName"`
	RoutingKey   string `json:"routingKey"`
	Headers      amqp.Table
}

type PublishTaskRequestData struct {
	Data interface{} `json:"data"`
}

func (req PublishTaskRequest) PublishMessage(ctx context.Context, amqpChannel *amqp.Channel) error {
	reqBytes := req.ReqBytes
	exchangeName := req.ExchangeName
	routingKey := req.RoutingKey

	if exchangeName == domain.EmptyString {
		log.Println("exchange_name_missing")
		return errors.New("exchange name missing")
	}
	if routingKey == domain.EmptyString {
		log.Println("routing_key_missing")
		return errors.New("routing key missing")
	}
	if reqBytes == nil {
		log.Println("No_Data_To_Publish")
		return errors.New("no data to publish")
	}

	err := amqpChannel.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Printf("Error: %v,\n failed_to_declare_a_rabbitMQ_exchange\n\n", err.Error())
		return errors.New("unable to declare exchange")
	}

	err = amqpChannel.Publish(
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         reqBytes,
		})
	if err != nil {
		log.Printf("Error: %v\n, failed_to_publish_a_rabbitMQ_message\n\n", err.Error())
		return errors.New("unable to publish messsage")
	}
	return nil
}
