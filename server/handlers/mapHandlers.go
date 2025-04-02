package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func ReverseGeocode(client *maps.Client, w http.ResponseWriter, r *http.Request, lat string, long string) {
	lat64, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		http.Error(w, `{"error": "Failed to convert latitude to float64"}`, http.StatusInternalServerError)
		return
	}
	long64, err := strconv.ParseFloat(long, 64)
	if err != nil {
		http.Error(w, `{"error": "Failed to convert longitude to float64"}`, http.StatusInternalServerError)
		return
	}
	result, err := mapping.ReverseGeocode(client, lat64, long64)
	if err != nil {
		http.Error(w, `{"error": "Failed to reverse geocode coordinates"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Reverse Geocoding successful!",
		"route":   result,
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
