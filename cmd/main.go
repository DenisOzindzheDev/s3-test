package main

import (
	"log"
	restapi "s3-test/internal/api"
	"s3-test/internal/application"
	"s3-test/internal/db/postgres"
	"s3-test/internal/s3/provider"
)

func main() {
	//creds minio
	//endpoint := "http://localhost:9000"
	filestore, err := provider.NewMinioProvider(
		"localhost:9000",
		"root",
		"password",
		false,
	)
	if err != nil {
		log.Fatal(err)
	}
	db, cls, err := postgres.NewDataBase("user=postgres password=postgres dbname=postgres host=localhost port=5432 sslmode=disable")
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
