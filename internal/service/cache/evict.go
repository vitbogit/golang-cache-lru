package cache

import (
	"context"

	"github.com/rs/zerolog/log"
)

// Evict обеспечивает ручное удаление данных по ключу
func (s *service) Evict(ctx context.Context, key string) (value interface{}, err error) {
	value, err = s.cacheRepository.Evict(ctx, key)
	if err != nil {
		log.Error().Err(err).Msg("ошибка удаления записи из кэша")
		return nil, err
	}

	return value, nil
}
