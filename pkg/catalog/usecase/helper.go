package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"message_queue_system/rabbitmq"
	"message_queue_system/rabbitmq/publisher"

	"github.com/streadway/amqp"
)

func PublishProductIdToQueue(ctx context.Context, conn *amqp.Connection, productId int) error {
	defer conn.Close()
	if conn.IsClosed() {
		log.Printf("Closed Connect")

		err := rabbitmq.Connect()
		if err != nil {
			log.Printf("Error: %v, unable to init rabbitmq", err.Error())
			return errors.New("unable to init rabbitmq conn")
		}
		conn = rabbitmq.Conn
	}
	amqpChannel, err := conn.Channel()
	if err != nil {
		log.Printf("Error: %v,\n failed_to_create_channel", err.Error())
		return errors.New("unable to create channel")
	}
	defer amqpChannel.Close()

	publishData := publisher.PublishTaskRequestData{}
	publishData.Data = productId
	reqBytes, err := json.Marshal(publishData)
	if err != nil {
		log.Printf("Error: %v,\n invalid_json_format", err.Error())
		return errors.New("invalid json format")
	}

	publishRequest := publisher.PublishTaskRequest{}
	publishRequest.QueueName = "store_product_images"
	publishRequest.ExchangeName = "store_product"
	publishRequest.RoutingKey = "store_product_images"
	publishRequest.ReqBytes = reqBytes
	err = publishRequest.PublishMessage(ctx, amqpChannel)
	if err != nil {
		log.Printf("Error: %v,\n failed_to_publish_message\n\n", err.Error())
		return errors.New("unable to publish message")
	}
	return nil
}
