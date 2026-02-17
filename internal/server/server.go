package server

import (
	"log"
	"morse_decoder_vanocry/internal/handlers"
	"net/http"
	"time"
)

// Server структура сервера с логгером и http-сервером
type Server struct {
	Logger *log.Logger
	Server *http.Server
}

// Новый экземпляр сервера
func NewServer(logger *log.Logger) *Server {
	// Роутер
	mux := http.NewServeMux()
	// Раздача статических файлов из папки ../cmd (иначе не грузит upload.html)
	// Static -> cmd ибо тесты не работают со static
	// mux.Handle("/cmd/", http.StripPrefix("/cmd/", http.FileServer(http.Dir("../cmd")))

	// Регистрация хендлеров
	mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("/upload", handlers.UploadHandler)

	// Экземпляр http.Server с настройками
	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{
		Logger: logger,
		Server: httpServer,
	}
}

// Run функция запуска сервера
func (s *Server) Run() error {
	s.Logger.Printf("Сервер запущен на порту 8080")
	return s.Server.ListenAndServe()
}

// Shutdown функция остановки сервера
func (s *Server) Shutdown() error {
	s.Logger.Printf("Остановка сервера...")
	return s.Server.Close()
}
