package app

import (
	"github.com/vitbogit/golang-cache-lru/internal/api/cache"
	"github.com/vitbogit/golang-cache-lru/internal/config"
	"github.com/vitbogit/golang-cache-lru/internal/repository"
	cacheRepository "github.com/vitbogit/golang-cache-lru/internal/repository/cache"
	"github.com/vitbogit/golang-cache-lru/internal/service"
	cacheService "github.com/vitbogit/golang-cache-lru/internal/service/cache"
)

// serviceProvider используется для корректного подтягивания зависимых частей приложения.
// Ожидается, что поля serviceProvider будут запрашиваться только через методы структуры,
// а сами методы будут внутри себя описывать, наличие каких других подтянутых зависимых частей необходимо,
// чтобы вернуть их самих.
type serviceProvider struct {
	httpConfig  config.HTTPConfig  // Конфиг HTTP-сервера
	cacheConfig config.CacheConfig // Конфиг кэша
	appConfig   config.AppConfig   // Конфиг приложения (общие настройки)

	cacheRepository repository.ILRUCache // Кэш база данных

	cacheService service.CacheService // Сервисный слой приложения

	cacheImpl *cache.Implementation // Имплементация API
}

// newServiceProvider создает пустой serviceProvider
func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// HTTPConfig возвращает конфиг http-сервера, предварительно проверив его наличие и наличие всех
// связанных с ним зависимых частей приложения, а в случае отсутствия чего-либо осуществляет
// попытку дозагрузки.
func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg := config.NewHTTPConfig()

		s.httpConfig = cfg
	}

	return s.httpConfig
}

// CacheConfig возвращает конфиг БД (кэша), предварительно проверив его наличие и наличие всех
// связанных с ним зависимых частей приложения, а в случае отсутствия чего-либо осуществляет
// попытку дозагрузки.
func (s *serviceProvider) CacheConfig() config.CacheConfig {
	if s.cacheConfig == nil {
		cfg := config.NewCacheConfig()

		s.cacheConfig = cfg
	}

	return s.cacheConfig
}

// AppConfig возвращает конфиг приложения (общие настройки), предварительно проверив его наличие и наличие всех
// связанных с ним зависимых частей приложения, а в случае отсутствия чего-либо осуществляет
// попытку дозагрузки.
func (s *serviceProvider) AppConfig() config.AppConfig {
	if s.appConfig == nil {
		cfg := config.NewAppConfig()

		s.appConfig = cfg
	}

	return s.appConfig
}

// CacheRepository возвращает БД (кэщ), предварительно проверив ее наличие и наличие всех
// связанных с ним зависимых частей приложения, а в случае отсутствия чего-либо осуществляет
// попытку дозагрузки.
func (s *serviceProvider) CacheRepository() repository.ILRUCache {
	if s.cacheRepository == nil {
		s.cacheRepository = cacheRepository.NewCache(s.CacheConfig().Size(), s.CacheConfig().DefaultTTL())
	}

	return s.cacheRepository
}

// CacheService возвращает сервис приложения, предварительно проверив его наличие и наличие всех
// связанных с ним зависимых частей приложения, а в случае отсутствия чего-либо осуществляет
// попытку дозагрузки.
func (s *serviceProvider) CacheService() service.CacheService {
	if s.cacheService == nil {
		s.cacheService = cacheService.NewService(
			s.CacheRepository(),
		)
	}

	return s.cacheService
}

// CacheImpl возвращает имплементацию API приложения, предварительно проверив ее наличие и наличие всех
// связанных с ним зависимых частей приложения, а в случае отсутствия чего-либо осуществляет
// попытку дозагрузки.
func (s *serviceProvider) CacheImpl() *cache.Implementation {
	if s.cacheImpl == nil {
		s.cacheImpl = cache.NewImplementation(s.CacheService())
	}

	return s.cacheImpl
}
