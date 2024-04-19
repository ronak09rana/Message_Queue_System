package controller

import (
	"context"

	"github.com/streadway/amqp"
)

type IProductController interface {
	ProcessProductImages(ctx context.Context, data interface{}, msg amqp.Delivery)
}
