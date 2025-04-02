package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"upgraded-telegram/main.go/server/services"
	"upgraded-telegram/main.go/server/services/mapping"

	"googlemaps.github.io/maps"
)

func GetDirections(client *maps.Client, w http.ResponseWriter, r *http.Request, a string, b string) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, `{"error": "Authorization header"}`, http.StatusInternalServerError)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		http.Error(w, `{"error": "Invalid token format"}`, http.StatusInternalServerError)
		return
	}
	claims := services.ParseAccessToken(token)
	if claims == nil {
		http.Error(w, `{"error": "Failed to verify token"}`, http.StatusInternalServerError)
		return
	}

	route, err := mapping.GetRoute(client, a, b)

	response := map[string]interface{}{
		"message": "Route Found!",
		"route":   route,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}
