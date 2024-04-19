package bootstrap

import (
	"message_queue_system/db"
	"message_queue_system/domain/interfaces"
	"message_queue_system/domain/interfaces/controller"
	"message_queue_system/domain/interfaces/repository"
	"message_queue_system/domain/interfaces/usecase"
	_productController "message_queue_system/pkg/catalog/controller"
	_productRepo "message_queue_system/pkg/catalog/repository"
	_productRoutes "message_queue_system/pkg/catalog/routes"
	_productUCase "message_queue_system/pkg/catalog/usecase"
	_userRepo "message_queue_system/pkg/user/repository"
	_userUCase "message_queue_system/pkg/user/usecase"
	"message_queue_system/rabbitmq"
	_consumer "message_queue_system/rabbitmq/consumer"

	"github.com/gin-gonic/gin"
)

var (
	userRepo          repository.IUserRepo
	userUCase         usecase.IUserUCase
	productRepo       repository.IProductRepo
	productUCase      usecase.IProductUCase
	productController controller.IProductController
	consumer          interfaces.IConsumer
)

func initRepos() {
	userRepo = _userRepo.NewUserRepository(db.Client)
	productRepo = _productRepo.NewProductRepository(db.Client)
}

func initControllers() {
	productController = _productController.NewProductController(productRepo)
}

func initConsumer() {
	consumer = _consumer.NewConsumerLayer(productController)
}

func initUCase() {
	userUCase = _userUCase.NewUserUCase(userRepo)
	productUCase = _productUCase.NewProductUCase(rabbitmq.Conn, userUCase, productRepo)
}

func initAPIs(apiGroup *gin.RouterGroup) {
	_productRoutes.InitRoutes(apiGroup, productUCase)
}

func Init(apiGroup *gin.RouterGroup) {
	initRepos()
	initControllers()
	initUCase()
	initAPIs(apiGroup)

	initConsumer()
	consumer.StartConsumers()
}
