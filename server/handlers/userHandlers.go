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
	"github.com/golang-jwt/jwt"
)

func CreateUser(client *dynamodb.Client, w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user db.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	id, err := uuid.NewV1()
	if err != nil {
		http.Error(w, `{"error": "Error generating user id"}`, http.StatusInternalServerError)
		return
	}

	email := strings.ToLower(user.Email)

	userId := fmt.Sprintf("u_%s", id)

	hashedPassword, err := services.HashedPassword(user.Password)
	if err != nil {
		return
	}

	newUser := map[string]types.AttributeValue{
		"id":       &types.AttributeValueMemberS{Value: userId},
		"name":     &types.AttributeValueMemberS{Value: user.Name},
		"email":    &types.AttributeValueMemberS{Value: email},
		"password": &types.AttributeValueMemberS{Value: hashedPassword},
	}

	err = db.CreateUser(client, "users", newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf(`{"message": "User created successfully", "user.id": %s}`, userId)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(message))

}

func AuthUser(client *dynamodb.Client, w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := db.GetUserByEmail(client, "users", req.Email)
	if err != nil {
		http.Error(w, `{"error": "No user found with that email."}`, http.StatusInternalServerError)
	}

	pass := services.CheckPasswordHash(req.Password, user.Password)
	if !pass {
		http.Error(w, `{"error": "Password Verfication Failed"}`, http.StatusInternalServerError)
		return
	}

	userClaims := services.UserClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	token, err := services.NewAccessToken(userClaims)
	if err != nil {
		http.Error(w, `{"error": "Password Verfication Failed"}`, http.StatusInternalServerError)
		return
	}

	refreshToken, err := services.NewRefreshToken(userClaims.StandardClaims)
	if err != nil {
		http.Error(w, `{"error": "Password Verfication Failed"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":       "Login Success",
		"token":         token,
		"refresh_token": refreshToken,
		"user":          user,
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

func GetAllUsers(client *dynamodb.Client, w http.ResponseWriter, r *http.Request) {

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
	resp, err := db.GetAllUsers(client, "users")
	if err != nil {
		http.Error(w, `{"error": "Failed to get all users"}`, http.StatusInternalServerError)
		return
	}
	var users []db.User

	for _, item := range resp {
		var user db.User
		err = attributevalue.UnmarshalMap(item, &user)
		if err != nil {
			http.Error(w, `{"error": "Failed to decode users"}`, http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	response := map[string]interface{}{
		"message": "Users Found!",
		"users":   users,
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

func GetUserByID(client *dynamodb.Client, w http.ResponseWriter, r *http.Request, id string) {

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
	resp, err := db.GetUserById(client, "users", id)
	if err != nil {
		http.Error(w, `{"error": "Failed to get all users"}`, http.StatusInternalServerError)
		return
	}
	var user db.User
	err = attributevalue.UnmarshalMap(resp, &user)
	if err != nil {
		http.Error(w, `{"error": "Failed to decode user"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User Found!",
		"user":    user,
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

func UpdateUser(client *dynamodb.Client, w http.ResponseWriter, r *http.Request) {

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

	var user db.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := db.UpdateUser(client, "users", user)
	if err != nil {
		http.Error(w, `{"error": "Failed to update user"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User Updated!",
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

func UpdatePassword(client *dynamodb.Client, w http.ResponseWriter, r *http.Request) {

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

	var user db.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	hashedPassword, err := services.HashedPassword(user.Password)
	if err != nil {
		return
	}

	user.Password = hashedPassword

	err = db.UpdatePassword(client, "users", user)
	if err != nil {
		http.Error(w, `{"error": "Failed to update user password"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User Password Updated!",
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

func DeleteUser(client *dynamodb.Client, w http.ResponseWriter, r *http.Request, id string) {

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

	err := db.DeleteUser(client, "users", id)
	if err != nil {
		http.Error(w, `{"error": "Failed to update user"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User Updated!",
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
