package usecase

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockIUserUCase struct {
	mock.Mock
}

func (m MockIUserUCase) FetchUser(ctx context.Context, userId int) (bool, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(bool), args.Error(1)
}
