package usecase

import "context"

type IUserUCase interface {
	FetchUser(ctx context.Context, userId int) (bool, error)
}
