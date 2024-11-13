package home

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func HelloHandler(w http.ResponseWriter, _ *http.Request) {
	// Charge les variables d'environnement depuis le fichier .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erreur lors du chargement du fichier .env : %v", err)
	}

	dbURL := os.Getenv("KEY_TEST")
	fmt.Fprintf(w, dbURL)
}

func MyFirstEndpoint(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "This is my first endpoint -->")
}
