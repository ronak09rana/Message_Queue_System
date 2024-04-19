package controller

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	createFile = CreateFile
)

func DownloadAndSaveImage(ctx context.Context, imagesArr []string) ([]string, error) {
	localImagesPathArr := make([]string, 0)

	outputDir := "./images"
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Printf("Error: %v,\n failed_to_make_directory", err.Error())
		return nil, errors.New("unable to create directory")
	}

	for _, imageUrl := range imagesArr {
		fileName := filepath.Base(imageUrl)
		outputPath := filepath.Join(outputDir, fileName)

		err := createFile(imageUrl, outputPath)
		if err != nil {
			log.Printf("Error: %v,\n failed_to_create_local_file", err.Error())
			return nil, errors.New("unable to create local file")
		}
		localImagesPathArr = append(localImagesPathArr, outputPath)
	}
	return localImagesPathArr, nil
}

func CreateFile(originalUrl, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	response, err := http.Get(originalUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}
