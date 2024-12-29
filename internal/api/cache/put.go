package cache

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vitbogit/golang-cache-lru/internal/converter"
	desc "github.com/vitbogit/golang-cache-lru/pkg/cache_v1"
)

// Put обеспечивает запись данных в кэш
func (i *Implementation) Put(w http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()
	log.Debug().Msg("API implementation method Put() requested by: " + r.Method + " " + r.URL.Path)
	defer func() {
		log.Debug().Msg("API implementation method Put() done with time " + time.Since(timeStart).String())
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var rawData desc.EntryPutData
	err = json.Unmarshal(body, &rawData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	convertedData := converter.ToEntryPutDataFromDesc(rawData)

	if len(convertedData.Key) == 0 || convertedData.TTL < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = i.cacheService.Put(context.Background(), convertedData.Key, convertedData.Value, convertedData.TTL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
