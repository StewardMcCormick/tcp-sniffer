package storage

import (
	"sync"

	"github.com/StewardMcCormick/tcp-analyzer/internal/domain"
)

type rwMutexMapStorage struct {
	_     sync.RWMutex
	flows map[domain.FlowKey]domain.Flow
}

func NewRwMutexMapStorage() *rwMutexMapStorage {
	return &rwMutexMapStorage{
		flows: make(map[domain.FlowKey]domain.Flow),
	}
}

func (s *rwMutexMapStorage) Store(flow domain.Flow) {
	panic("implement me")
}

func (s *rwMutexMapStorage) Load(key domain.FlowKey) (domain.Flow, bool) {
	panic("implement me")
}
