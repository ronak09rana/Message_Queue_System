package usecase

import (
	"context"
	"errors"
	"log"
	"message_queue_system/domain/interfaces/repository"
	"message_queue_system/domain/interfaces/usecase"
)

type UserUCase struct {
	UserRepo repository.IUserRepo
}

func NewUserUCase(userRepo repository.IUserRepo) usecase.IUserUCase {
	return UserUCase{
		UserRepo: userRepo,
	}
}

func (uuc UserUCase) FetchUser(ctx context.Context, userId int) (bool, error) {
	userExists, err := uuc.UserRepo.Get(ctx, userId)
	if err != nil {
		log.Printf("Error: %v, unable_to_fetch_user\n\n", err.Error())
		return false, errors.New("unable to fetch user")
	}
	return userExists, nil
}
