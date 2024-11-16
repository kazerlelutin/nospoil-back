package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func TmdbAPIHandler(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erreur lors du chargement du fichier .env : %v", err)
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		http.Error(w, "Clé API non définie", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	endpoint := vars["endpoint"]

	apiURL := fmt.Sprintf("https://api.themoviedb.org/3/%s", endpoint)

	queryParams := r.URL.Query()
	if len(queryParams) > 0 {
		apiURL += "?" + queryParams.Encode()
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la création de la requête : %v", err), http.StatusInternalServerError)
		return
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'appel API : %v", err), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Erreur de l'API : statut %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la lecture de la réponse : %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(body))
}
