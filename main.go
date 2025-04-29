package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Requête reçue: %s %s depuis %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r) // Passer la requête au handler suivant ou au handler final
	})
}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Bienvenue sur l'API publique !"))
}

func main() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de la base de données: %v", err)
	}
	defer db.Close()
	router := mux.NewRouter()

	router.Use(loggingMiddleware)
	router.HandleFunc("/", HomeHandler).Methods("GET")

	port := ":8080"
	log.Printf("Le serveur écoute sur le port %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
