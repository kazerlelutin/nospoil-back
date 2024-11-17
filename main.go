package main

import (
	"k-space-go/handlers"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Middleware pour les CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*") //TODO WHITE LIST
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		//OPTIONS
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()

	r.Use(enableCORS)

	r.HandleFunc("/", handlers.HelloHandler).Methods("GET")
	r.HandleFunc("/tv/{endpoint:.*}", handlers.TmdbAPIHandler).Methods("GET")

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
