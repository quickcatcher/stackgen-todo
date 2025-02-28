package main

import (
	"context"
	"fmt"
	"stackgen-todo/core/controller"
	"stackgen-todo/routes"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	endpoint := "127.0.0.1:9000"
	accessKeyID := "admin"
	secretAccessKey := "password"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		fmt.Println("Error while connecting to minio ", err)
	}
	bucketName := "newbucket"
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		fmt.Println("Error while looking for bucket ", err)
		return
	}
	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			fmt.Println("Error while creating bucket")
			return
		}
		fmt.Println("Bucket created succesfully")
	}

	todoList := controller.NewTodoList(bucketName, minioClient)

	r := gin.Default()
	router := r.Group("/stackgen")
	routes.EngineRoutes(router, todoList)

	r.Run(":8000")

}
