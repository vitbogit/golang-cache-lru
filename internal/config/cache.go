package config

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	cacheSizeEnvName        = "CACHE_SIZE"        // Имя переменной окружения для параметра размера кэша
	cacheSizeFlagName       = "cache-size"        // Имя флага для параметра размера кэша
	cacheDefaultTTLEnvName  = "DEFAULT_CACHE_TTL" // Имя переменной окружения для параметра TTL по умолчанию кэша
	cacheDefaultTTLFlagName = "default-cache-ttl" // Имя флага для параметра TTL по умолчанию кэша
)

// CacheConfig описывает методы конфига кэша
type CacheConfig interface {
	Size() int                 // Размер кэша
	DefaultTTL() time.Duration // TTL по умолчанию
}

// cacheConfig задает поля конфига кэша
type cacheConfig struct {
	size       int           // Размер кэша
	defaultTTL time.Duration // TTL по умолчанию
}

// cacheConfigJSON задает поля конфига кэша, описанные в JSON (ограниченный набор типов)
type cacheConfigJSON struct {
	Size       int    `json:"cache_size"`         // Размер кэша
	DefaultTTL string ` json:"default_cache_ttl"` // TTL по умолчанию (строка)
}

// CacheDefaultValues загружает значения по умолчанию для кэша из JSON-файла
func CacheDefaultValues() cacheConfigJSON {
	defaultValuesFile, err := LoadJSON(cacheCfgDefaultValuesPath)
	if err != nil {
		log.Fatal().Err(err).Msg("ошибка при чтении файла конфигурации кэша со значениями по умолчанию")
	}

	var defaultValues cacheConfigJSON
	err = json.Unmarshal(defaultValuesFile, &defaultValues)
	if err != nil {
		log.Fatal().Err(err).Msg("ошибка при чтении файла конфигурации кэша со значениями по умолчанию")
	}

	return defaultValues
}

// NewCacheConfig собирает актуальный конфиг кэша по трехступенчатому принципу
//
// - Если для параметра определен флаг запуска, используется он
//
// - Если флаг не определен, используется переменная окружения
//
// - Если не определены ни флаг, ни переменная окружения, используется значение по умолчанию
func NewCacheConfig() CacheConfig {
	var err error

	// flag value
	sizeFlag := flags.cacheSize
	defaultTTLFlag := flags.cacheDefaultTTL

	// env value
	sizeEnv := os.Getenv(cacheSizeEnvName)
	defaultTTLEnv := os.Getenv(cacheDefaultTTLEnvName)

	// default values
	defaultValues := CacheDefaultValues()

	// Трехступенчатый выбор размера кэша
	var size int
	switch {
	case sizeFlag != 0:
		size = sizeFlag
	case len(sizeEnv) > 0:
		size, err = strconv.Atoi(sizeEnv)
		if err != nil {
			log.Fatal().Msg("некорректный формат размер кэша для приложения (считан из переменной среды)")
		}
	case defaultValues.Size != 0:
		size = defaultValues.Size
	default:
		log.Fatal().Msg("не удалось определить значение параметра размер кэша для приложения")
	}

	if size <= 0 {
		log.Fatal().Msg("некорректный формат размер кэша для приложения, size должен быть > 0")
	}

	// Трехступенчатый выбор TTL по умолчанию
	var defaultTTLString string
	switch {
	case len(defaultTTLFlag) > 0:
		defaultTTLString = defaultTTLFlag
	case len(defaultTTLEnv) > 0:
		defaultTTLString = defaultTTLEnv
	case len(defaultValues.DefaultTTL) > 0:
		defaultTTLString = defaultValues.DefaultTTL
	default:
		log.Fatal().Msg("не удалось определить значение параметра TTL по умолчанию в кэше для приложения")
	}
	defaultTTL, err := time.ParseDuration(defaultTTLString)
	if err != nil {
		log.Fatal().Msg("некорректный формат параметра TTL по умолчанию в кэше, должен являться временем")
	}

	if defaultTTL <= 0 {
		log.Fatal().Msg("некорректный формат параметра TTL по умолчанию в кэше, default TTL должен быть > 0")
	}

	return &cacheConfig{
		size:       size,
		defaultTTL: defaultTTL,
	}
}

// Size возвращает параметр размер кэша из конфига
func (cfg *cacheConfig) Size() int {
	return cfg.size
}

// DefaultTTL возвращает параметр TTL по умолчанию из конфига
func (cfg *cacheConfig) DefaultTTL() time.Duration {
	return cfg.defaultTTL
}
