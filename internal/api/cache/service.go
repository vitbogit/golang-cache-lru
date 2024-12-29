// Package cache содержит определение имплементации API сервиса
package cache

import (
	"github.com/vitbogit/golang-cache-lru/internal/service"
)

// Implementation задает поля в имплементации API сервиса
type Implementation struct {
	cacheService service.CacheService
}

// NewImplementation создает новую имплементацию
func NewImplementation(cacheService service.CacheService) *Implementation {
	return &Implementation{
		cacheService: cacheService,
	}
}
