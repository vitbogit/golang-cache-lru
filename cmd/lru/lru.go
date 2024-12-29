// Package main содержит входную точку в приложение
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/vitbogit/golang-cache-lru/internal/app"
)

func main() {

	// Создаём контекст, который завершится при получении сигналов SIGINT/SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Канал для получения сигналов
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем обработчик сигналов в отдельной горутине
	go func() {
		sig := <-signalChan
		log.Info().Str("signal", sig.String()).Msg("получили termination signal")
		cancel() // Завершаем контекст
	}()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("не удалось инициализировать приложение")
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("не удалось запустить приложение")
	}

	<-ctx.Done()
	log.Info().Msg("совершен gracefull shutdown")
}
