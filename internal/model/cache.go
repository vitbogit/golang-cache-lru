// Package model содержит определения Entities слоя приложения.
// В виду простоты приложения и факта, что оно не будет являться полноценным проектом,
// некоторые промежуточные структуры и конвертеры для них могут быть опущены.
package model

import "time"

// EntryPutData представляет поля для записи значения в кэш на уровне Entities
type EntryPutData struct {
	Key   string
	Value interface{}
	TTL   time.Duration
}
