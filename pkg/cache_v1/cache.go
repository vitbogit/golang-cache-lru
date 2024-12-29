// Package cache предоставляет структуры, использующиеся сервисом golang-cache-lru в API-слое, то есть для взаимодействия с внешним миром.
package cache

// EntryPutData представляет набор полей, используемых для создания записи в кэше.
type EntryPutData struct {
	Key        string      `json:"key"`         // Ключ
	Value      interface{} `json:"value"`       // Значение
	TTLSeconds int         `json:"ttl_seconds"` // TTL (в секундах)
}

// EntryGetAllData описывает результат запроса на получение всех ключей и их значений их кэша.
// Получение всего наполнения кэша происходит в виде двух слайсов: слайса ключей и слайса значений.
// Пары ключ-значения из кэша располагаются на соответствующих позициях в слайсах.
type EntryGetAllData struct {
	Keys   []string      `json:"keys"`
	Values []interface{} `json:"values"`
}

// EntryGetData представляет набор полей, которые сервис возвращает в качестве данных о записи в кэше.
type EntryGetData struct {
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
	ExpiresAt int64       `json:"expires_at"`
}
