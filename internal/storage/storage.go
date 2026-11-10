package storage

import "github.com/StewardMcCormick/tcp-analyzer/internal/domain"

type FlowStorage interface {
	Store(flow domain.Flow)
	Load(key domain.FlowKey) (domain.Flow, bool)
}
