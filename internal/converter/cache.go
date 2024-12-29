// Package converter содержит конвертеры между слоями приложения.
// В виду простоты приложения и факта, что оно не будет являться полноценным проектом,
// некоторые промежуточные структуры и конвертеры для них могут быть опущены.
package converter

import (
	"time"

	"github.com/vitbogit/golang-cache-lru/internal/model"
	desc "github.com/vitbogit/golang-cache-lru/pkg/cache_v1"
)

// ToEntryPutDataFromDesc конвертирует поля для создания новой записи в кэше из API-слоя в Entities
func ToEntryPutDataFromDesc(info desc.EntryPutData) model.EntryPutData {
	return model.EntryPutData{
		Key:   info.Key,
		Value: info.Value,
		TTL:   time.Second * time.Duration(info.TTLSeconds),
	}
}
