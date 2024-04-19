package consumer

import (
	"context"
	"encoding/json"
	"log"
	"message_queue_system/domain/entity"
	"message_queue_system/domain/interfaces"
	"message_queue_system/domain/interfaces/controller"
	"message_queue_system/rabbitmq"
	"message_queue_system/rabbitmq/publisher"
	"os"

	"github.com/streadway/amqp"
)

type ConsumerLayer struct {
	ProductController controller.IProductController
}

func NewConsumerLayer(productController controller.IProductController) interfaces.IConsumer {
	return ConsumerLayer{
		ProductController: productController,
	}
}

func (cl ConsumerLayer) StartConsumers() {
	consumers := []entity.Consumer{
		{
			QueueName:    "store_product_images",
			ExchangeName: "store_product",
			ExchangeType: "direct",
			RoutingKey:   "store_product_images",
		},
	}

	for _, consumer := range consumers {
		go cl.StartConsumer(consumer)
	}
}

func (cl ConsumerLayer) StartConsumer(consumer entity.Consumer) {
	conn := rabbitmq.Conn
	defer conn.Close()
	exchangeName := consumer.ExchangeName
	exchangeType := consumer.ExchangeType
	amqpChannel, err := conn.Channel()
	if err != nil {
		log.Printf("Error: %v,\n unable to connect to rabbitmq channel", err.Error())
		os.Exit(0)
	}
	defer amqpChannel.Close()

	err = amqpChannel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Printf("Error: %v, unable to declare exchange", err.Error())
		os.Exit(0)
	}

	queue, err := amqpChannel.QueueDeclare(
		consumer.QueueName, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Printf("Error: %v,\n unable to declare queue", err.Error())
		os.Exit(0)
	}

	err = amqpChannel.QueueBind(
		consumer.QueueName,    // queue name
		consumer.RoutingKey,   // routing key
		consumer.ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Printf("Error: %v,\n unable to bind queue", err.Error())
		os.Exit(0)
	}

	msgs, err := amqpChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Printf("Error: %v, unable to consume message", err.Error())
		os.Exit(0)
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received a message for Queue: %v", consumer.QueueName)
			//time.Sleep(30 * time.Second)
			cl.ConsumeMessage(queue.Name, msg, consumer)
		}
	}()
}

func (cl ConsumerLayer) ConsumeMessage(queueName string, msg amqp.Delivery, consumer entity.Consumer) {
	ctx := context.Background()
	reqBytes := msg.Body
	var consumedData publisher.PublishTaskRequestData
	err := json.Unmarshal(reqBytes, &consumedData)
	if err != nil {
		log.Printf("Error: %v,\n unable to unmarshal consumed message", err.Error())
		msg.Ack(false)
		return
	}

	switch queueName {
	case "store_product_images":
		go cl.ProductController.ProcessProductImages(ctx, consumedData.Data, msg)
	default:
		return
	}
}
