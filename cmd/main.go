package main

import (
	"TaskManager/internal/api"
	"TaskManager/internal/config"
	"TaskManager/internal/service"
	"TaskManager/internal/storage"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DataBaseURL())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	store := storage.NewPostgresTaskStorage(pool)
	svc := service.NewTaskService(store)
	handler := api.NewTaskHandler(svc)

	http.Handle("/tasks", handler)

	serverPort := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Println("listening on ", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}
