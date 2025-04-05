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

func CreateOrder(client *dynamodb.Client, w http.ResponseWriter, r *http.Request) {

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

	var order db.Order
	err := json.NewDecoder(r.Body).Decode(&order)
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

	orderId := fmt.Sprintf("o_%s", id)

	if order.User == "" {
		http.Error(w, `{"error": "No user for this order found."}`, http.StatusInternalServerError)
		return
	}

	if order.Status == "" {
		http.Error(w, `{"error": "No status for this order found."}`, http.StatusInternalServerError)
		return
	}

	newOrder := map[string]types.AttributeValue{
		"id":         &types.AttributeValueMemberS{Value: orderId},
		"created_at": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().UnixMilli())},
	}

	if len(order.Items) > 0 {
		newOrder["items"] = &types.AttributeValueMemberSS{Value: order.Items}
	}

	err = db.CreateOrder(client, "items", newOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":  "Order craeted!",
		"order.id": orderId,
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

func GetOrderById(client *dynamodb.Client, w http.ResponseWriter, r *http.Request, id string) {

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

	resp, err := db.GetItemById(client, "items", id)
	if err != nil {
		http.Error(w, `{"error": "Failed to get order"}`, http.StatusInternalServerError)
		return
	}

	var order db.Order
	err = attributevalue.UnmarshalMap(resp, &order)
	if err != nil {
		http.Error(w, `{"error": "Failed to decode order"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Got order!",
		"order":   order,
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

func GetAllOrders(client *dynamodb.Client, w http.ResponseWriter, r *http.Request) {

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

	resp, err := db.GetAllOrders(client, "orders")
	if err != nil {
		http.Error(w, `{"error": "Failed to get all orders"}`, http.StatusInternalServerError)
		return
	}

	var orders []db.Order
	err = attributevalue.UnmarshalListOfMaps(resp, &orders)
	if err != nil {
		http.Error(w, `{"error": "Failed to decode orders"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Got orders!",
		"orders":  orders,
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

func UpdateOrder(client *dynamodb.Client, w http.ResponseWriter, r *http.Request) {

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

	var order db.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = db.UpdateOrders(client, "orders", order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := `{"message": "Order updated"}`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func DeleteOrder(client *dynamodb.Client, w http.ResponseWriter, r *http.Request, id string) {

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

	err := db.DeleteOrder(client, "orders", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := `{"message": "Order deleted"}`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
