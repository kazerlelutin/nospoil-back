package main

import (
	"k-space-go/handlers"
	utils "k-space-go/utils"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	utils.Db()

	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HelloHandler).Methods("GET")
	r.HandleFunc("/tv/{endpoint:.*}", handlers.TmdbAPIHandler).Methods("GET")

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":8080", r))

}
