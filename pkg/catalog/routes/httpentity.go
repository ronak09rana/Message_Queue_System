package routes

import (
	"errors"
	"message_queue_system/domain/entity"
)

type CreateProductReq struct {
	UserId      int      `json:"userId"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	Price       int      `json:"price"`
}

func (req CreateProductReq) validate() error {
	if req.UserId == 0 || len(req.Name) == 0 || req.Price == 0 {
		return errors.New("mandatory details missing")
	}
	return nil
}

func (req CreateProductReq) toProductDto() entity.Product {
	return entity.Product{
		UserId:      req.UserId,
		Name:        req.Name,
		Description: req.Description,
		Images:      req.Images,
		Price:       req.Price,
	}
}
