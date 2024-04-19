package usecase

import (
	"context"
	"message_queue_system/domain/entity"
	mockProductRepo "message_queue_system/mocks/catalog/repository"
	mockUserUCase "message_queue_system/mocks/user/usecase"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

func Benchmark_UpsertProduct(b *testing.B) {
	ctx := context.Background()
	product := entity.Product{
		UserId:      1,
		Name:        "Pizza",
		Description: "Spanish Delicious Dish",
		Images: []string{
			"https://unsplash.com/photos/mVnft46UU1Y",
			"https://unsplash.com/photos/TfOxBg75im4",
		},
		Price: 200,
	}

	userUCaseMock := &mockUserUCase.MockIUserUCase{}
	productRepoMock := &mockProductRepo.MockIProductRepo{}
	conn := &amqp.Connection{}

	userUCaseMock.On("FetchUser", mock.Anything, mock.Anything).
		Return(true, nil)

	productRepoMock.On("Upsert", mock.Anything, mock.Anything).
		Return(1, nil)

	publishProductIdToQueue = func(ctx context.Context, conn *amqp.Connection, productId int) error {
		return nil
	}

	puc := ProductUCase{
		UserUCase:   userUCaseMock,
		ProductRepo: productRepoMock,
		Conn:        conn,
	}

	for i := 0; i < b.N; i++ {
		puc.UpsertProduct(ctx, product)
	}
}
