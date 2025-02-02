package main

import (
	"fmt"
	"log"
	"net/http"

	_ "music/docs"
	"music/internal/controller"
	"music/internal/db"
	"music/internal/repository"
	"music/internal/services"
	"music/pkg/config"
	"music/pkg/logger"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Music API
// @version 1.0
// @description API for managing music library
// @host localhost:8080
// @BasePath /
func main() {
	// Инициализация конфигурации
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	fmt.Printf("All config: %+v\n", cfg)

	// Инициализация логгера
	_log := logger.NewLogger()

	// Подключение к базе данных
	dbConn, err := db.NewDB(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer dbConn.Close()

	// Инициализация репозитория
	repo := repository.NewMainRepository(dbConn.PostgreSQL)

	// Инициализация сервиса
	service := services.NewMainService(repo, http.DefaultClient, cfg, _log)

	mux := mux.NewRouter()

	// Инициализация контроллера
	ctrl := controller.NewMainController(service, mux, _log)

	ctrl.RegisterHandlers()

	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/swagger.json")
	})
	
	mux.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"), 
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Запуск сервера
	fmt.Printf("The Music API is running: http://localhost:%s\n\tPress Ctl + C for stopping\n", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, mux); err != nil {
		log.Print(err.Error())
	}
}
