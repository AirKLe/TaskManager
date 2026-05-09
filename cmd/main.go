package main

import (
	"TaskManager/internal/api"
	"TaskManager/internal/service"
	"TaskManager/internal/storage"
	"log"
	"net/http"
)

func main() {
	store := storage.NewInMemoryTaskStorage()
	svc := service.NewTaskService(store)
	handler := api.NewTaskHandler(svc)

	http.Handle("/tasks", handler)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
