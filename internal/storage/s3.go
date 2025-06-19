package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	Client     *s3.Client
	Bucket     string
	Uploader   *manager.Uploader
	BaseURL    string
	FolderPath string
}

func NewS3Storage(client *s3.Client, bucket, baseURL, folder string) *S3Storage {
	return &S3Storage{
		Client:     client,
		Bucket:     bucket,
		Uploader:   manager.NewUploader(client),
		BaseURL:    baseURL, // например, https://s3.amazonaws.com/your-bucket
		FolderPath: folder,
	}
}

func (s *S3Storage) UploadFile(file multipart.File, filename string) (string, error) {
	key := filepath.Join(s.FolderPath, fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename))

	_, err := s.Uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   file,
		ACL:    "public-read",
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s", s.BaseURL, key)
	return url, nil
}

func (s *S3Storage) DeleteFile(path string) error {
	key := filepath.Base(path) // может потребоваться кастомная логика
	_, err := s.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(filepath.Join(s.FolderPath, key)),
	})
	return err
}
