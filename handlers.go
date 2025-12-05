package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (app *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"status": "auth-service ok"}
	json.NewEncoder(w).Encode(response)
}

func (app *App) validateKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}

	if _, exists := app.Keys[key]; !exists {
		http.Error(w, "invalid key", http.StatusUnauthorized)
		return
	}

	response := map[string]string{"message": "Key valid"}
	json.NewEncoder(w).Encode(response)
}

func (app *App) createKeyHandler(w http.ResponseWriter, r *http.Request) {
	key, err := generateKey()
	if err != nil {
		http.Error(w, "cannot generate key", http.StatusInternalServerError)
		return
	}

	app.Keys[key] = KeyInfo{Key: key}
	response := map[string]string{"key": key}
	json.NewEncoder(w).Encode(response)
}

func (app *App) masterKeyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		masterKey := r.Header.Get("X-Master-Key")

		if masterKey == "" || masterKey != app.MasterKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
