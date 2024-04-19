package controller

import (
	"context"
	mockProductRepo "message_queue_system/mocks/catalog/repository"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

func Benchmark_ProcessProductImages(b *testing.B) {
	defer unMockBenchmarkFunction()
	ctx := context.Background()
	data := 1
	msg := &amqp.Delivery{}

	productRepoMock := &mockProductRepo.MockIProductRepo{}

	productRepoMock.On("Get", mock.Anything, mock.Anything).
		Return([]string{"image1.jpg", "image2.jpg"}, nil)

	productRepoMock.On("Save", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	downloadAndSaveImage = func(ctx context.Context, imagesArr []string) ([]string, error) {
		return []string{"imageLocal1.png", "image2Local.png"}, nil
	}

	pc := ProductController{
		ProductRepo: productRepoMock,
	}

	for i := 0; i < b.N; i++ {
		pc.ProcessProductImages(ctx, data, *msg)
	}
}

func Test_DownloadAndSaveImage(t *testing.T) {
	defer unMockUnitTestFunction()
	ctx := context.Background()
	imagesInputArr := []string{
		"https://unsplash.com/photos/mVnft46UU1Y",
		"https://unsplash.com/photos/TfOxBg75im4",
	}

	createFile = func(inputPath, outputPath string) error {
		return nil
	}

	DownloadAndSaveImage(ctx, imagesInputArr)
}

func unMockBenchmarkFunction() {
	downloadAndSaveImage = DownloadAndSaveImage
}

func unMockUnitTestFunction() {
	createFile = CreateFile
}
