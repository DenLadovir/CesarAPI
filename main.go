package main

import (
	"CesarAPI/database"
	"CesarAPI/middleware"
	"CesarAPI/routes"
	"context"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	if err := database.ConnectDatabase(); err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v\n", err)
	}

	r := chi.NewRouter()

	r.Group(func(c chi.Router) {
		c.Use(middleware.JWTMiddleware)
		taskHandler := &routes.TaskHandler{}
		taskHandler.RegisterRoutes(r, database.DB)
	})

	// Открытые маршруты (без middleware)
	routes.AuthRoutes(r)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	go func() {
		log.Infof("Запуск сервера на порту %s...", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка при запуске сервера: %v\n", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении работы сервера: %v\n", err)
	}
	log.Info("Сервер корректно завершил работу.")
}
