package mapping

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"googlemaps.github.io/maps"
)

// repo and documentation https://github.com/googlemaps/google-maps-services-go?tab=readme-ov-file

func FindMaps() *maps.Client {
	err := godotenv.Load("server/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mapsKey := os.Getenv("GOOGLE_MAPS_API_KEY")

	mapClient, err := maps.NewClient(maps.WithAPIKey(mapsKey))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	log.Println("Connecting to Google Maps")
	return mapClient
}

func GetRoute(client *maps.Client, origin string, destination string) ([]maps.Route, error) {
	req := &maps.DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Perth",
	}
	route, _, err := client.Directions(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return route, nil
}

func ReverseGeocode(client *maps.Client, lat float64, long float64) ([]maps.GeocodingResult, error) {
	r := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{Lat: lat, Lng: long},
	}
	reverseGeocodingResponse, err := client.ReverseGeocode(context.Background(), r)
	if err != nil {
		return nil, err
	}
	return reverseGeocodingResponse, nil
}
