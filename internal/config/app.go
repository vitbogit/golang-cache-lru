package config

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	AppLogLevelEnvName  = "LOG_LEVEL" // Имя переменной окружения для параметра уровень логов приложения
	AppLogLevelFlagName = "log-level" // Имя флага для параметра уровень логов приложения
)

// AppConfig описывает методы конфига приложения (общие настройки)
type AppConfig interface {
	LogLevel() string // Уровень логирования (например, "WARN")
}

// appConfig задает поля конфига приложения (общие настройки)
type appConfig struct {
	logLevel string // Уровень логирования (например, "WARN")
}

// appConfigJSON задает поля конфига приложения (общие настройки), описанные в JSON (ограниченный набор типов)
type appConfigJSON struct {
	LogLevel string `json:"log_level"` // Уровень логирования (например, "WARN")
}

// AppDefaultValues загружает значения по умолчанию для общих настроек приложения из JSON-файла
func AppDefaultValues() appConfigJSON {
	defaultValuesFile, err := LoadJSON(appCfgDefaultValuesPath)
	if err != nil {
		log.Fatal().Err(err).Msg("ошибка при чтении файла конфигурации общих настроек приложения со значениями по умолчанию")
	}

	var defaultValues appConfigJSON
	err = json.Unmarshal(defaultValuesFile, &defaultValues)
	if err != nil {
		log.Fatal().Err(err).Msg("ошибка при чтении файла конфигурации общих настроек приложения со значениями по умолчанию")
	}

	return defaultValues
}

// # NewAppConfig собирает актуальный конфиг общих настроек приложения по трехступенчатому принципу
//
// - Если для параметра определен флаг запуска, используется он
//
// - Если флаг не определен, используется переменная окружения
//
// - Если не определены ни флаг, ни переменная окружения, используется значение по умолчанию
//
// # Важное замечание про log level
//
// В выбранной архитектуре приложения пришлось прибегнуть к считыванию log level еще до вызова текущей функции NewAppConfig(),
// чтобы корректно вывести логи на ранних этапах запуска приложения.
//
// Это значит, что повторно выбранное значение log level в NewAppConfig() уже не будет использоваться для установки уровня логирования,
// но для формальности его запишем в структуру конфига сервера.
//
// На всякий случай осуществляется проверка повторно считанного log level на равенство УЖЕ УСТАНОВЛЕННОМУ уровню логирования.
func NewAppConfig() AppConfig {

	// flag value
	logLevelFlag := flags.logLevel

	// env value
	logLevelEnv := os.Getenv(AppLogLevelEnvName)

	// default values
	defaultValues := AppDefaultValues()

	//  Трехступенчатый выбор
	var logLevel string
	switch {
	case len(logLevelFlag) > 0:
		logLevel = logLevelFlag
	case len(logLevelEnv) > 0:
		logLevel = logLevelEnv
	case len(defaultValues.LogLevel) > 0:
		logLevel = defaultValues.LogLevel
	default:
		log.Fatal().Msg("не удалось определить значение уровня логирования для конфига приложения (*повторное считывание)")
	}

	if providedLvl, err := zerolog.ParseLevel(logLevel); err != nil {
		log.Fatal().Msg("некорректный формат значения уровня логирования для конфига приложения (*повторное считывание)")
	} else if providedLvl != zerolog.GlobalLevel() {
		log.Fatal().Msg("повторно считанный уровень логирования для конфига приложения не совпал с установленным на раннем этапе")
	}

	return &appConfig{
		logLevel: logLevel,
	}
}

// LogLevel возвращает параметр уровень логирования приложения из конфига
func (cfg *appConfig) LogLevel() string {
	return cfg.logLevel
}
