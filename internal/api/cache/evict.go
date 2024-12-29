package cache

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

// Evict обеспечивает ручное удаление данных по ключу
func (i *Implementation) Evict(w http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()
	log.Debug().Msg("API implementation method Evict() requested by: " + r.Method + " " + r.URL.Path)
	defer func() {
		log.Debug().Msg("API implementation method Evict() done with time " + time.Since(timeStart).String())
	}()

	key := chi.URLParam(r, "key")
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := i.cacheService.Evict(context.Background(), key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if value == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
