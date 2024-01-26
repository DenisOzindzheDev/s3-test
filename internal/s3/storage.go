package s3

import (
	"context"
	"s3-test/internal/models"
)

type ImageStorage interface {
	Connect() error
	UploadFile(context.Context, models.ImageUnit) (string, error)
	DownloadFile(context.Context, string) (models.ImageUnit, error)
}
