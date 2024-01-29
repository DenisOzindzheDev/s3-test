package main

import (
	"log"
	"os"
	restapi "s3-test/internal/api"
	"s3-test/internal/application"
	"s3-test/internal/db/postgres"
	"s3-test/internal/s3/provider"
)

func main() {
	filestore, err := provider.NewMinioProvider(
		os.Getenv("MINIO_HOST"),
		os.Getenv("MINIO_USER"),
		os.Getenv("MINIO_ACCES_KEY"),
		false,
	)
	if err != nil {
		log.Fatal(err)
	}
	db, cls, err := postgres.NewDataBase(os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatal(err)
	}
	defer cls()

	core, err := application.NewSkater(filestore, db, "http://localhost:9090/browser/testbucket")
	if err != nil {
		log.Fatal(err)
	}
	srv := restapi.NewServer(core)

	err = srv.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}

}
