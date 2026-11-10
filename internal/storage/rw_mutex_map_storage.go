package storage

import (
	"sync"

	"github.com/StewardMcCormick/tcp-analyzer/internal/domain"
)

type rwMutexMapStorage struct {
	mu    sync.RWMutex
	flows map[domain.FlowKey]domain.Flow
}

func NewRwMutexMapStorage() *rwMutexMapStorage {
	return &rwMutexMapStorage{
		flows: make(map[domain.FlowKey]domain.Flow),
	}
}

func (s *rwMutexMapStorage) Store(flow domain.Flow) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.flows[flow.Key] = flow
}

func (s *rwMutexMapStorage) Load(key domain.FlowKey) (domain.Flow, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	flow, ok := s.flows[key]
	return flow, ok
}
