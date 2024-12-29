// Package config содержит структуры, методы и константы для конфигурирования приложения.
//
// # Сборка конфигов осуществляется по трехступенчатому принципу
//
// - Если для параметра определен флаг запуска, используется он
//
// - Если флаг не определен, используется переменная окружения
//
// - Если не определены ни флаг, ни переменная окружения, используется значение по умолчанию
package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

const (
	appCfgDefaultValuesPath   = "configs/app.json"   // Путь к значения по умолчанию для настроек приложения (общих настроек)
	cacheCfgDefaultValuesPath = "configs/cache.json" // Путь к значения по умолчанию для настроек непосредственно кэша
	httpCfgDefaultValuesPath  = "configs/http.json"  // Путь к значения по умолчанию для настроек непосредственно сервера приложения
	cfgEnvPath                = ".env"               // Путь к конфигурационному файлу среды
)

// LoadEnv оборачивает функцию Load из godotenv, которая читает конфигурационный файл среды, обработкой ошибок
func LoadEnv() error {
	err := godotenv.Load(cfgEnvPath)
	if err != nil {
		log.Warn().Err(fmt.Errorf("error loading .env: %v", err))
		return err
	}

	return nil
}

// flags хранит считанные флаги
//
// TODO: реализовать трехступенчатую логику логирования без использования глобальной переменной.
var flags Flags

// Flags описывает флаги приложения
type Flags struct {
	cacheSize       int    // Размер кэша
	cacheDefaultTTL string // TTL по умолчанию в кэше

	httpHostPort string // Хост-порт HTTP-сервера

	logLevel string // Уровень логирования
}

// GetLogFlag возвращает флаг с уровнем логирования
func GetLogFlag() string {
	return flags.logLevel
}

// LoadFlags загружает флаги приложения в общую структуру
func LoadFlags() {
	size := flag.Int(cacheSizeFlagName, 0, "an int")
	defaultTTL := flag.String(cacheDefaultTTLFlagName, "", "a string")

	hostPort := flag.String(httpHostPortFlagName, "", "a string")

	logLevel := flag.String(AppLogLevelFlagName, "", "a string")

	flag.Parse()

	flags = Flags{
		cacheSize:       *size,
		cacheDefaultTTL: *defaultTTL,
		httpHostPort:    *hostPort,
		logLevel:        *logLevel,
	}
}

// LoadJSON считывает JSON-файл конфигурации по указанному пути, эта функция вызывается отдельно для каждого файла.
//
// TODO: сделать более абстрактную функцию, чтобы в результате можно было получать не последовательность байтов,
// а сами структуры, например.
func LoadJSON(path string) (fileBytes []byte, err error) {
	fileBytes, err = os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
