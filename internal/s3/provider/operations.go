package provider

import (
	"context"
	"log"
	"s3-test/internal/application"
	"s3-test/internal/models"

	"github.com/minio/minio-go/v7"
)

func (s *MinioProvider) UploadFile(ctx context.Context, file models.ImageUnit) (string, error) {
	filename := application.GenerateObjectName(file.User)
	_, err := s.client.PutObject(
		ctx,
		UserObjectBucketName,
		filename,
		file.Payload,
		file.PayloadSize,
		minio.PutObjectOptions{ContentType: "image/png"},
	)
	return filename, err
}
func (s *MinioProvider) DownloadFile(ctx context.Context, file string) (models.ImageUnit, error) {
	reader, err := s.client.GetObject(
		ctx,
		UserObjectBucketName,
		file,
		minio.GetObjectOptions{},
	)
	if err != nil {
		log.Fatalf("Error downloading image: %v", err)
	}
	defer reader.Close()

	return models.ImageUnit{}, nil
}
