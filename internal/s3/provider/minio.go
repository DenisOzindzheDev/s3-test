package provider

import (
	"log"
	"s3-test/internal/s3"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioProvider struct {
	minioAuthData
	client *minio.Client
}

type minioAuthData struct {
	url      string
	user     string
	password string
	token    string
	ssl      bool
}

func NewMinioProvider(minioURL string, minioUser string, minioPassword string, ssl bool) (s3.ImageStorage, error) {
	return &MinioProvider{
		minioAuthData: minioAuthData{
			password: minioPassword,
			url:      minioURL,
			user:     minioUser,
			ssl:      ssl,
		}}, nil

}

func (s *MinioProvider) Connect() error {
	var err error
	s.client, err = minio.New(s.url, &minio.Options{
		Creds:  credentials.NewStaticV4(s.user, s.password, ""),
		Secure: s.ssl,
	})
	if err != nil {
		log.Fatal(err)
	}
	return err
}
