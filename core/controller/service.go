package controller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/cast"
)

func CreateToDoItem(item *Item, todoList *Todo) (res *CreateTodoItemResponse, err error) {
	todoList.MutexLock()
	defer todoList.MutexUnlock()

	todoList.LastID++
	item.ID = todoList.LastID
	todoList.Todos[todoList.LastID] = item

	res = &CreateTodoItemResponse{
		Id: todoList.LastID,
	}
	return
}

func GetTodoItem(id int, todoList *Todo) (item *Item, err error) {
	todoList.MutexLock()
	defer todoList.MutexUnlock()

	item, ok := todoList.Todos[id]
	if !ok {
		return nil, nil
	}
	reqParams := make(map[string][]string)
	for i, attachment := range item.Attachments {
		presignedURL, e := todoList.MinioClient.PresignedGetObject(context.Background(), todoList.BucketName, attachment, time.Second*60*60, reqParams)
		if err != nil {
			return nil, e
		}
		fmt.Println(presignedURL)
		item.Attachments[i] = presignedURL.String()
	}
	return
}

func GetAllItems(todoList *Todo) (res *GetAllTodoItemsResponse, err error) {
	if todoList.LastID == 0 {
		return
	}
	id := 1
	var item *Item
	res = &GetAllTodoItemsResponse{}
	for id < todoList.LastID {
		item, err = GetTodoItem(id, todoList)
		if err != nil {
			return
		}
		res.Items = append(res.Items, item)
	}
	return
}

func AddAttachment(id int, todoList *Todo, files []*multipart.FileHeader) (err error) {
	todoList.MutexLock()
	defer todoList.MutexUnlock()

	for _, file := range files {
		src, e := file.Open()
		if e != nil {
			return e
		}
		defer src.Close()

		fileContents, e := io.ReadAll(src)
		if e != nil {
			return e
		}

		objectName := cast.ToString(id) + "_" + file.Filename

		reader := bytes.NewReader(fileContents)
		_, err = todoList.MinioClient.PutObject(context.Background(), todoList.BucketName, objectName, reader, int64(len(fileContents)), minio.PutObjectOptions{
			ContentType: file.Header["Content-Type"][0],
		})
		if err != nil {
			return
		}
		todoList.Todos[id].Attachments = append(todoList.Todos[id].Attachments, objectName)
	}
	return
}

func UpdateTodoItem(id int, item *Item, todoList *Todo) {
	todoList.MutexLock()
	defer todoList.MutexUnlock()

	if _, ok := todoList.Todos[id]; !ok {
		return
	}
	todoList.Todos[id].Title = item.Title
	todoList.Todos[id].Description = item.Description
}
