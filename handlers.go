package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// healthHandler verifica se o serviço está operacional
func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "auth-service ok"})
}

// validateKeyHandler valida uma API key enviada via Authorization: Bearer <key>
func (a *App) validateKeyHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	keyString := strings.TrimPrefix(authHeader, "Bearer ")

	if keyString == "" {
		http.Error(w, "Authorization header ausente", http.StatusUnauthorized)
		return
	}

	// Hash da chave
	keyHash := hashAPIKey(keyString)

	// Verifica no banco se existe uma chave ativa com esse hash
	var id int
	err := a.DB.QueryRow(
		"SELECT id FROM api_keys WHERE key_hash = $1 AND is_active = true",
		keyHash,
	).Scan(&id)

	if err != nil {
		log.Printf("Chave inválida (hash %s): %v", keyHash[:8], err)
		http.Error(w, "Chave inválida ou inativa", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Chave válida"})
}

// createKeyHandler cria uma nova API key
func (a *App) createKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Corpo inválido", http.StatusBadRequest)
		return
	}

	if body.Name == "" {
		http.Error(w, "'name' é obrigatório", http.StatusBadRequest)
		return
	}

	// Gera chave e hash
	keyPlain, err := generateAPIKey()
	if err != nil {
		http.Error(w, "Falha ao gerar chave", http.StatusInternalServerError)
		return
	}

	keyHash := hashAPIKey(keyPlain)

	var id int
	err = a.DB.QueryRow(
		"INSERT INTO api_keys (name, key_hash) VALUES ($1, $2) RETURNING id",
		body.Name,
		keyHash,
	).Scan(&id)

	if err != nil {
		log.Printf("Erro ao salvar chave: %v", err)
		http.Error(w, "Erro ao salvar chave", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"key":     keyPlain, // retorna a chave apenas uma vez!
		"name":    body.Name,
		"message": "Guarde esta chave com segura
