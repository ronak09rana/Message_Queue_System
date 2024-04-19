package repository

import "context"

type IUserRepo interface {
	Get(ctx context.Context, userId int) (bool, error)
}
