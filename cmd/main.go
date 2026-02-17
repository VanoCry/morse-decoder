package main

import (
	"log"
	"morse_decoder_vanocry/internal/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Логгер
	logger := log.New(os.Stdout, "SERVER: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Сервер
	srv := server.NewServer(logger)

	// Канал для обработки сигналов
	// Используется для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Запускаем сервер в горутине (всё также для graceful shutdown)
	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Ошибка при запуске сервера: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-quit
	logger.Println("Получен сигнал завершения, останавливаем сервер...")

	// Остановка сервера
	if err := srv.Shutdown(); err != nil {
		logger.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	logger.Println("Сервер успешно остановлен")

}
