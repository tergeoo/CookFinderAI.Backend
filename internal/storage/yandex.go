package storage

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log/slog"
	"mime/multipart"
	"time"
)

type YandexStorage struct {
	client     *minio.Client
	bucketName string
	endpoint   string
}

func NewYandexStorage(endpoint, accessKey, secretKey, bucket string) (*YandexStorage, error) {
	slog.Info("Creating Yandex Storage client", "endpoint", endpoint, "bucket", bucket)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}
	return &YandexStorage{
		client:     client,
		bucketName: bucket,
		endpoint:   endpoint,
	}, nil
}

func (it *YandexStorage) UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	objectName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)

	// Upload
	_, err := it.client.PutObject(ctx, it.bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: fileHeader.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", err
	}

	// Public URL
	url := fmt.Sprintf("https://%s/%s/%s", it.endpoint, it.bucketName, objectName)
	return url, nil
}

func (it *YandexStorage) DeleteFile(ctx context.Context, objectName string) error {
	err := it.client.RemoveObject(ctx, it.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from Yandex Cloud: %w", err)
	}
	return nil
}
