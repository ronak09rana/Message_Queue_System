package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"message_queue_system/domain/interfaces/repository"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) repository.IUserRepo {
	return UserRepo{
		DB: db,
	}
}

func (ur UserRepo) Get(ctx context.Context, userId int) (bool, error) {
	conn, err := ur.DB.Conn(ctx)
	if err != nil {
		log.Printf("Error: %v\n, unable_to_db_connect\n\n", err.Error())
		return false, errors.New("unable to db connect")
	}
	defer conn.Close()

	var name string
	sqlQuery := "SELECT name FROM user WHERE id = ?"
	args := []interface{}{userId}
	err = conn.QueryRowContext(ctx, sqlQuery, args...).Scan(&name)
	if err != nil {
		log.Printf("Error: %v\n, unable_to_fetch_user\n\n", err)
		return false, errors.New("unable to fetch user")
	}
	return true, nil
}
