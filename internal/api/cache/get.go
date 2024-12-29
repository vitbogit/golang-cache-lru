package cache

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	desc "github.com/vitbogit/golang-cache-lru/pkg/cache_v1"
)

// Get обеспечивает получение данных из кэша по ключу
func (i *Implementation) Get(w http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()
	log.Debug().Msg("API implementation method Get() requested by: " + r.Method + " " + r.URL.Path)
	defer func() {
		log.Debug().Msg("API implementation method Get() done with time " + time.Since(timeStart).String())
	}()

	key := chi.URLParam(r, "key")
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, expiresAt, err := i.cacheService.Get(context.Background(), key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if value == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sendData := desc.EntryGetData{
		Key:       key,
		Value:     value,
		ExpiresAt: expiresAt.Unix(),
	}
	sendDataBytes, err := json.Marshal(sendData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(sendDataBytes)
	w.WriteHeader(http.StatusOK)
}
