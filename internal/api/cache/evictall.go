package cache

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// EvictAll  обеспечивает ручную инвалидацию всего кэша
func (i *Implementation) EvictAll(w http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()
	log.Debug().Msg("API implementation method EvictAll() requested by: " + r.Method + " " + r.URL.Path)
	defer func() {
		log.Debug().Msg("API implementation method EvictAll() done with time " + time.Since(timeStart).String())
	}()

	err := i.cacheService.EvictAll(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
