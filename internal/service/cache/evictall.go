package cache

import (
	"context"

	"github.com/rs/zerolog/log"
)

// EvictAll  обеспечивает ручную инвалидацию всего кэша
func (s *service) EvictAll(ctx context.Context) error {
	err := s.cacheRepository.EvictAll(ctx)
	if err != nil {
		log.Error().Err(err).Msg("ошибка очистки кэша")
		return err
	}

	return nil
}
