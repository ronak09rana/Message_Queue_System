package controller

import (
	"context"
	"encoding/json"
	"log"
	"message_queue_system/domain/interfaces/controller"
	"message_queue_system/domain/interfaces/repository"

	"github.com/streadway/amqp"
)

var (
	downloadAndSaveImage = DownloadAndSaveImage
)

type ProductController struct {
	ProductRepo repository.IProductRepo
}

func NewProductController(productRepo repository.IProductRepo) controller.IProductController {
	return ProductController{
		ProductRepo: productRepo,
	}
}

func (pc ProductController) ProcessProductImages(ctx context.Context, data interface{}, msg amqp.Delivery) {
	productIdBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error: %v,\n failed_to_marshal_request_body\n\n", err.Error())
		msg.Ack(false)
		return
	}

	var productId int
	err = json.Unmarshal(productIdBytes, &productId)
	if err != nil {
		log.Printf("Error: %v,\n failed_to_unmarshal_productIdBytes\n\n", err.Error())
		msg.Ack(false)
		return
	}

	productImagesArr, err := pc.ProductRepo.Get(ctx, productId)
	if err != nil {
		log.Printf("Error: %v,\n failed_to_fetch_product", err.Error())
		msg.Ack(false)
		return
	}

	localImagesPathArr, err := downloadAndSaveImage(ctx, productImagesArr)
	if err != nil {
		log.Printf("Error: %v,\n failed_to_download_and_save_image", err.Error())
		msg.Ack(false)
		return
	}

	err = pc.ProductRepo.Save(ctx, productId, localImagesPathArr)
	if err != nil {
		log.Printf("Error: %v,\n failed_to_save_images_local_path", err.Error())
		msg.Ack(false)
	}
}
