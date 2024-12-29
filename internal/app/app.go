// Package app содержит основные верхнеуровневые элементы приложения,
// включая определение структуры App (приложения), функций для запуска приложения и сервера, а также механизмы для корректного
// подтягивания зависимых частей, таких как БД (сам кэш) и конфиги.
package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vitbogit/golang-cache-lru/internal/config"
)

const (
	serverShutdownTimeOut = 10 * time.Second
)

// App задает структуру приложения
type App struct {
	serviceProvider *serviceProvider // Менеджер зависимых частей приложения
	httpServer      *http.Server     // HTTP-сервер
}

// NewApp создает новое приложение и вызывает функцию для инициализации
func NewApp(ctx context.Context) (*App, error) {
	log.Debug().Msg("Creating new app")
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	log.Debug().Msg("Sucessfully created new app")
	return a, nil
}

// Run запускает приложение
func (a *App) Run(ctx context.Context) error {
	return a.runHTTPServer(ctx)
}

// initDeps инициализирует приложение
func (a *App) initDeps(ctx context.Context) error {
	log.Debug().Msg("Initing deps")

	inits := []func(context.Context) error{
		a.initConfigAndLogger,
		a.initServiceProvider,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	log.Debug().Msg("Sucessfully inited deps")
	return nil
}

// initConfigAndLogger подготавливает конфиги и логгер приложения
func (a *App) initConfigAndLogger(_ context.Context) error {
	log.Debug().Msg("Initing config and logger")

	// Подготовка конфигов из флагов
	config.LoadFlags()
	logLevelString := config.GetLogFlag() // уровень логирования из флагов

	// Подготовка конфигов из переменных среды
	errReadingEnv := config.LoadEnv()

	// Трехступенчатая логика выбора значения log_level из конфига,
	// до вызова общей инициализации всех остальных параметров (чтобы уже
	// сейчас установить уровень логирования)
	if len(logLevelString) == 0 {
		if errReadingEnv == nil {
			logLevelString = os.Getenv(config.AppLogLevelEnvName)
		}
		if len(logLevelString) == 0 {
			logLevelString = config.AppDefaultValues().LogLevel
			if len(logLevelString) == 0 {
				zerolog.SetGlobalLevel(zerolog.InfoLevel) // чтобы точно показались следующие логи
				log.Fatal().Msg("error reading log level")
			}
		}
	}

	// Преобразование строки в zerolog.Level
	level, errParsingLvl := zerolog.ParseLevel(logLevelString)
	if errParsingLvl != nil {
		log.Fatal().Msg("invalid log level provided")
	}

	zerolog.SetGlobalLevel(level)

	// Ранее не логированная ошибка, выводим после установки уровня конфига
	if errReadingEnv != nil {
		log.Warn().Err(fmt.Errorf("can`t parse .env: %v", errParsingLvl))
	}

	log.Debug().Msg("Sucessfully inited config and logger")
	return nil
}

// initServiceProvider инициализирует service provider
func (a *App) initServiceProvider(_ context.Context) error {
	log.Debug().Msg("Initing service provider")

	a.serviceProvider = newServiceProvider()

	// При использовании .Dict() в zerolog для красивого вывода JSON поле time выводилось в другом месте строки
	// и это не получилось быстро исправить, поэтому используется просто Sprintf
	log.Debug().Msg(fmt.Sprintf("using App config: %+v", a.serviceProvider.AppConfig()))
	log.Debug().Msg(fmt.Sprintf("using HTTP config: %+v", a.serviceProvider.HTTPConfig()))
	log.Debug().Msg(fmt.Sprintf("using Cache config: %+v", a.serviceProvider.CacheConfig()))

	log.Debug().Msg("Sucessfully inited service provider")
	return nil
}

// initServiceProvider инициализирует initHTTPServer
func (a *App) initHTTPServer(_ context.Context) error {
	log.Debug().Msg("Initing http server")

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("cache service homepage!"))
	})

	r.Route("/api/lru", func(r chi.Router) {
		r.Post("/", a.serviceProvider.CacheImpl().Put)

		r.Get("/{key}", a.serviceProvider.CacheImpl().Get)
		r.Get("/", a.serviceProvider.CacheImpl().GetAll)

		r.Delete("/{key}", a.serviceProvider.CacheImpl().Evict)
		r.Delete("/", a.serviceProvider.CacheImpl().EvictAll)
	})

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().HostPort(),
		Handler: r,
	}

	log.Debug().Msg("Sucessfully inited http server")
	return nil
}

// runHTTPServer запускает HTTP-сервер
func (a *App) runHTTPServer(ctx context.Context) error {
	// Запуск сервера в горутине
	go func() {
		log.Info().Msg(fmt.Sprintf("запуск HTTP сервера на %s", a.httpServer.Addr))
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("не удалось запустить сервер")
		}
	}()

	// Завершаем сервер при отмене контекста
	<-ctx.Done()
	log.Info().Msg("отключение HTTP сервера...")

	// Создаём таймаут для завершения сервера
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), serverShutdownTimeOut)
	defer shutdownCancel()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("server forced to shutdown")
	} else {
		log.Info().Msg("server gracefully stopped")
	}

	return nil
}
