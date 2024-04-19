package repository

import (
	"context"
	"message_queue_system/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockIProductRepo struct {
	mock.Mock
}

func (m MockIProductRepo) Upsert(ctx context.Context, product entity.Product) (int, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(int), args.Error(1)
}

func (m MockIProductRepo) Get(ctx context.Context, productId int) ([]string, error) {
	args := m.Called(ctx, productId)
	return args.Get(0).([]string), args.Error(1)
}

func (m MockIProductRepo) Save(ctx context.Context, productId int, paths []string) error {
	args := m.Called(ctx, productId, paths)
	return args.Error(0)
}
