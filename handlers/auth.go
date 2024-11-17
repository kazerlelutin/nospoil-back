package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func SignInWithOTP(w http.ResponseWriter, r *http.Request) {

	godotenv.Load()

	var requestBody struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Requête invalide, impossible de décoder le corps", http.StatusBadRequest)
		return
	}

	if requestBody.Email == "" {
		http.Error(w, "Email et code OTP sont requis", http.StatusBadRequest)
		return
	}

	email := requestBody.Email

	if email == "" {
		http.Error(w, "Email manquant dans la requête", http.StatusBadRequest)
		return
	}

	supabaseURL := os.Getenv("SUPABASE_SERVICE_URL")
	apiKey := os.Getenv("SUPABASE_SERVICE_KEY")

	if supabaseURL == "" || apiKey == "" {
		http.Error(w, "Paramètres d'environnement manquants", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s/auth/v1/otp", supabaseURL)

	payload := map[string]interface{}{
		"email":            email,
		"shouldCreateUser": true,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la sérialisation du payload : %v", err), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la création de la requête : %v", err), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		http.Error(w, fmt.Sprintf("Erreur lors de l'appel à l'API : %v", err), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	fmt.Fprintf(w, "Réponse de l'API : %s", string(body))
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la lecture de la réponse : %v", err), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		http.Error(w, fmt.Sprintf("Erreur de l'API : statut %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Un code OTP a été envoyé à %s", email)
}

func VerifyOTP(w http.ResponseWriter, r *http.Request) {

	godotenv.Load()

	var requestBody struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Requête invalide, impossible de décoder le corps", http.StatusBadRequest)
		return
	}

	if requestBody.Email == "" || requestBody.OTP == "" {
		http.Error(w, "Email et code OTP sont requis", http.StatusBadRequest)
		return
	}

	supabaseURL := os.Getenv("SUPABASE_SERVICE_URL")
	apiKey := os.Getenv("SUPABASE_SERVICE_KEY")
	if supabaseURL == "" || apiKey == "" {
		http.Error(w, "Paramètres d'environnement manquants", http.StatusInternalServerError)
		return
	}

	// Endpoint de vérification de l'OTP
	url := fmt.Sprintf("%s/auth/v1/verify", supabaseURL)

	// Corps de la requête pour vérifier le code OTP
	payload := map[string]string{
		"email": requestBody.Email,
		"token": requestBody.OTP, // Le code OTP que l'utilisateur a reçu
		"type":  "email",         // Pour indiquer que c'est un OTP envoyé par e-mail
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la sérialisation du payload : %v", err), http.StatusInternalServerError)
		return
	}

	// Création de la requête HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la création de la requête : %v", err), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Envoi de la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'appel à l'API : %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Lecture de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la lecture de la réponse : %v", err), http.StatusInternalServerError)
		return
	}

	// Vérification du code de statut de la réponse
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Erreur de vérification de l'OTP : %s", string(body)), http.StatusInternalServerError)
		return
	}

	// Réponse de succès
	fmt.Fprintf(w, "OTP vérifié avec succès : %s", string(body))
}
