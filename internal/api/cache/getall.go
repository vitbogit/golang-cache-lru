package cache

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	desc "github.com/vitbogit/golang-cache-lru/pkg/cache_v1"
)

// GetAll обеспечивает получение всего наполнения кэша в виде двух слайсов: слайса ключей и слайса значений.
// Пары ключ-значения из кэша располагаются на соответствующих позициях в слайсах.
func (i *Implementation) GetAll(w http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()
	log.Debug().Msg("API implementation method GetAll() requested by: " + r.Method + " " + r.URL.Path)
	defer func() {
		log.Debug().Msg("API implementation method GetAll() done with time " + time.Since(timeStart).String())
	}()

	keys, values, err := i.cacheService.GetAll(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(keys) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	sendData := desc.EntryGetAllData{
		Keys:   keys,
		Values: values,
	}
	sendDataBytes, err := json.Marshal(sendData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(sendDataBytes)
}
