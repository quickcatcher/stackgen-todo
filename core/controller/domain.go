package controller

import (
	"sync"

	"github.com/minio/minio-go/v7"
)

type Item struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Attachments []string `json:"attachments"`
}

type Todo struct {
	Todos       map[int]*Item
	LastID      int
	mu          sync.Mutex
	BucketName  string
	MinioClient *minio.Client
}

func NewTodoList(bucketName string, m *minio.Client) *Todo {
	return &Todo{
		Todos:       make(map[int]*Item),
		LastID:      0,
		BucketName:  bucketName,
		MinioClient: m,
	}
}

func (s *Todo) MutexLock() {
	s.mu.Lock()
}

func (s *Todo) MutexUnlock() {
	s.mu.Unlock()
}
