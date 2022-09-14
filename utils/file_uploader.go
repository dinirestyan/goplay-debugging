package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

var path string = "assets/"

func UploadFile(c *gin.Context, folder, name string, file *multipart.FileHeader) (string, error) {
	// check id folder doesn't exist
	if _, err := os.Stat(path + folder); os.IsNotExist(err) {
		if err := os.MkdirAll(path+folder, os.ModePerm); err != nil {
			return "", errors.New("can't create folder with error " + err.Error())
		}
	}

	filename := filepath.Base(file.Filename)
	ext := filepath.Ext(filename)

	if !(ext == ".jpeg" || ext == ".jpg" || ext == ".png" || ext == ".mp4") {
		return "", errors.New("You can only upload image and video files")
	}

	if err := c.SaveUploadedFile(file, path+folder+"/"+filename); err != nil {
		return "", errors.New("can't upload file with error " + err.Error())
	}

	url, err := GoogleUploadObject(path+folder+"/"+filename, filename)
	if err != nil {
		return "", err
	}
	if err := os.Remove(path + folder + "/" + filename); err != nil {
		return "", errors.New("failed to delete file : " + err.Error())
	}
	return url, nil
}

func GoogleUploadObject(path, fname string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", errors.New("Failed to create google cloud client: " + err.Error())
	}
	defer client.Close()

	f, err := os.Open(path)
	if err != nil {
		return "", errors.New("Failed to open file: " + err.Error())
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := client.Bucket(os.Getenv("GOOGLE_BUCKET")).Object(fname).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return "", errors.New("Error copying file to cloud: " + err.Error())
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}
	return fname, nil
}
