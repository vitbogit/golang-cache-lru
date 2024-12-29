// Package cache содержит определение структур и логику работы сервисного слоя приложения кэша golang-cahe-lru.
package cache

import (
	"github.com/vitbogit/golang-cache-lru/internal/repository"
	def "github.com/vitbogit/golang-cache-lru/internal/service"
)

var _ def.CacheService = (*service)(nil)

// service структура сервиса golang-cahe-lru
type service struct {
	cacheRepository repository.ILRUCache
}

// NewService создает новый сервис golang-cahe-lru
func NewService(
	cacheRepository repository.ILRUCache,
) *service {
	return &service{
		cacheRepository: cacheRepository,
	}
}
