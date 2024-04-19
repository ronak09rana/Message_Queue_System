package repository

import (
	"context"
	"message_queue_system/domain/entity"
)

type IProductRepo interface {
	Upsert(ctx context.Context, product entity.Product) (int, error)
	Get(ctx context.Context, productId int) ([]string, error)
	Save(ctx context.Context, productId int, imagesArr []string) error
}
