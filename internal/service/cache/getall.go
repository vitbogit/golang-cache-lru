package cache

import (
	"context"

	"github.com/rs/zerolog/log"
)

// GetAll обеспечивает получение всего наполнения кэша в виде двух слайсов: слайса ключей и слайса значений.
// Пары ключ-значения из кэша располагаются на соответствующих позициях в слайсах.
func (s *service) GetAll(ctx context.Context) (keys []string, values []interface{}, err error) {
	keys, values, err = s.cacheRepository.GetAll(ctx)
	if err != nil {
		log.Error().Err(err).Msg("ошибка получения всего наполнения кэша")
		return nil, nil, err
	}

	return keys, values, nil
}
