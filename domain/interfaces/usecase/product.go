package usecase

import (
	"context"
	"message_queue_system/domain/entity"
)

type IProductUCase interface {
	UpsertProduct(ctx context.Context, product entity.Product) error
}
