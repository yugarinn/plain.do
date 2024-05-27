package main

import (
	"html/template"
	"net/http"
	"log"
	"path/filepath"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseFiles(
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", "index.html"),
	))
}

func renderTemplate(w http.ResponseWriter, templateName string) {
	err := templates.ExecuteTemplate(w, templateName, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "base.html")
}

func main() {
	http.HandleFunc("/", indexHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
