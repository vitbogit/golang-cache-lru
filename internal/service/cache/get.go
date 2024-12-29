package cache

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

// Get обеспечивает получение данных из кэша по ключу
func (s *service) Get(ctx context.Context, key string) (value interface{}, expiresAt time.Time, err error) {
	value, expiresAt, err = s.cacheRepository.Get(ctx, key)
	if err != nil {
		log.Error().Err(err).Msg("ошибка получения записи из кэша")
		return nil, time.Time{}, err
	}

	return value, expiresAt, nil
}
