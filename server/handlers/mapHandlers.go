package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"googlemaps.github.io/maps"
)

func GetDirections(client *maps.Client, w http.ResponseWriter, r *http.Request) {
	req := &maps.DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Perth",
	}
	route, _, err := client.Directions(context.Background(), req)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

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
