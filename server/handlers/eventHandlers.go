package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"upgraded-telegram/main.go/server/services"
	"upgraded-telegram/main.go/server/services/db"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gofrs/uuid"
)

func CreateEventHandler(client *dynamodb.Client, w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
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

	var event db.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	id, err := uuid.NewV1()
	if err != nil {
		http.Error(w, `{"error": "Error generating user id"}`, http.StatusInternalServerError)
		return
	}

	eventId := fmt.Sprintf("e_%s", id)

	if event.Name == "" {
		http.Error(w, `{"error": "Event Name is required"}`, http.StatusInternalServerError)
		return
	}

	if event.StartDate == 0 {
		http.Error(w, `{"error": "Event Start Date is required"}`, http.StatusInternalServerError)
		return
	}

	if event.EndDate == 0 {
		http.Error(w, `{"error": "Event End Date is required"}`, http.StatusInternalServerError)
		return
	}

	newEvent := map[string]types.AttributeValue{
		"id":         &types.AttributeValueMemberS{Value: eventId},
		"name":       &types.AttributeValueMemberS{Value: event.Name},
		"start_date": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", event.StartDate)},
		"emd_date":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", event.EndDate)},
		"active":     &types.AttributeValueMemberBOOL{Value: false},
		"created_at": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().UnixMilli())},
	}

	if len(event.AssignedTo) > 0 {
		newEvent["assigned_to"] = &types.AttributeValueMemberSS{Value: event.AssignedTo}
	}

	if event.Description != "" {
		newEvent["description"] = &types.AttributeValueMemberS{Value: event.Description}
	}

	if event.LocationName != "" {
		newEvent["location_name"] = &types.AttributeValueMemberS{Value: event.LocationName}
	}

	if event.LocationAddress != "" {
		newEvent["location_address"] = &types.AttributeValueMemberS{Value: event.LocationAddress}
	}

	if event.LocationLat != 0 {
		newEvent["location_lat"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", event.LocationLat)}
	}

	if event.LocationLong != 0 {
		newEvent["location_long"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", event.LocationLong)}
	}

	if event.Notes != "" {
		newEvent["notes"] = &types.AttributeValueMemberS{Value: event.Notes}
	}

	if event.FirstNotification != 0 {
		newEvent["first_notification"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", event.FirstNotification)}
	}

	if event.SecondNotification != 0 {
		newEvent["second_notification"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", event.SecondNotification)}
	}

	err = db.CreateEvent(client, "events", newEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":  "Event saved!",
		"event.id": eventId,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

func GetEventById(client *dynamodb.Client, w http.ResponseWriter, r *http.Request, id string) {

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

	resp, err := db.GetEventById(client, "events", id)
	if err != nil {
		http.Error(w, `{"error": "Failed to get event"}`, http.StatusInternalServerError)
		return
	}

	var event db.Event
	err = attributevalue.UnmarshalMap(resp, &event)
	if err != nil {
		http.Error(w, `{"error": "Failed to decode event"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Got event!",
		"event":   event,
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

func GetAllEvents(client *dynamodb.Client, w http.ResponseWriter, r *http.Request) {

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

	resp, err := db.GetAllEvents(client, "events")
	if err != nil {
		http.Error(w, `{"error": "Failed to get all events"}`, http.StatusInternalServerError)
		return
	}

	var events []db.Event
	err = attributevalue.UnmarshalListOfMaps(resp, &events)
	if err != nil {
		http.Error(w, `{"error": "Failed to decode events"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Got evemts!",
		"events":  events,
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

func UpdateEvent(client *dynamodb.Client, w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
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

	var event db.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = db.UpdateEvent(client, "events", event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := `{"message": "Event updated"}`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func DeleteEvent(client *dynamodb.Client, w http.ResponseWriter, r *http.Request, id string) {

	if r.Method != http.MethodDelete {
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

	err := db.DeleteEvent(client, "events", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := `{"message": "Event deleted"}`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
