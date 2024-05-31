package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yugarinn/plain.do/handlers"
	"github.com/yugarinn/plain.do/repository"
)

func main() {
	repository.InitDB("plain.db")
	repository.Migrate()

	r := mux.NewRouter()

	r.HandleFunc("/", handlers.IndexHandler).Methods("GET")
	r.HandleFunc("/about", handlers.AboutHandler).Methods("GET")
	r.HandleFunc("/{listHash}", handlers.IndexHandler).Methods("GET")
	r.HandleFunc("/todos/{todoID}", handlers.DeleteTodoHandler).Methods("DELETE")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/ws", handlers.WsHandler)
	http.Handle("/", r)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

