package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

// Put обеспечивает запись данных в кэш
func (s *service) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if len(key) == 0 || ttl < 0 {
		log.Error().Msg("некорректные данные для добавления в кэш")
		return fmt.Errorf("некорректные входные данные")
	}

	err := s.cacheRepository.Put(ctx, key, value, ttl)
	if err != nil {
		log.Error().Err(err).Msg("ошибка добавления в кэш")
		return err
	}

	return nil
}
