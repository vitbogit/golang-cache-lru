package config

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

const (
	httpHostPortEnvName  = "SERVER_HOST_PORT" // Имя переменной окружения для параметра хост-порт сервера
	httpHostPortFlagName = "server-host-port" // Имя флага для параметра хост-порт сервера
)

// HTTPConfig описывает методы конфига сервера
type HTTPConfig interface {
	HostPort() string // Пара "хост:порт" одной строкой
}

// httpConfig задает поля конфига сервера
type httpConfig struct {
	hostPort string // Пара "хост:порт" одной строкой
}

// httpConfigJSON задает поля конфига сервера, описанные в JSON (ограниченный набор типов)
type httpConfigJSON struct {
	HostPort string `json:"server_host_port"` // Пара "хост:порт" одной строкой
}

// HTTPDefaultValues загружает значения по умолчанию для сервера приложения из JSON-файла
func HTTPDefaultValues() httpConfigJSON {
	defaultValuesFile, err := LoadJSON(httpCfgDefaultValuesPath)
	if err != nil {
		log.Fatal().Err(err).Msg("ошибка при чтении файла конфигурации сервера со значениями по умолчанию")
	}

	var defaultValues httpConfigJSON
	err = json.Unmarshal(defaultValuesFile, &defaultValues)
	if err != nil {
		log.Fatal().Err(err).Msg("ошибка при обработке файла конфигурации сервера со значениями по умолчанию")
	}

	return defaultValues
}

// NewHTTPConfig собирает актуальный конфиг сервера по трехступенчатому принципу
//
// - Если для параметра определен флаг запуска, используется он
//
// - Если флаг не определен, используется переменная окружения
//
// - Если не определены ни флаг, ни переменная окружения, используется значение по умолчанию
func NewHTTPConfig() HTTPConfig {
	// Flag value
	hostPortFlag := flags.httpHostPort

	// Env value
	hostPortEnv := os.Getenv(httpHostPortEnvName)

	// Default values
	defaultValues := HTTPDefaultValues()

	// Трехступенчатый выбор
	var hostPort string
	switch {
	case len(hostPortFlag) > 0:
		hostPort = hostPortFlag
	case len(hostPortEnv) > 0:
		hostPort = hostPortEnv
	case len(defaultValues.HostPort) > 0:
		hostPort = defaultValues.HostPort
	default:
		log.Fatal().Msg("не удалось определить значение параметра хост-порт для сервера приложения")
	}

	return &httpConfig{
		hostPort: hostPort,
	}
}

// HostPort возвращает хост-порт параметр настроек сервера
func (cfg *httpConfig) HostPort() string {
	return cfg.hostPort
}
