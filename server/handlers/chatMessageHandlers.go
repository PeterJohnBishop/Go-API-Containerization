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

func CreateChat(client *dynamodb.Client, w http.ResponseWriter, r *http.Request) {

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

	var chat db.Chat
	err := json.NewDecoder(r.Body).Decode(&chat)
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

	chatId := fmt.Sprintf("c_%s", id)

	if len(chat.Users) == 0 {
		http.Error(w, `{"error": "Users array must not be empty"}`, http.StatusInternalServerError)
		return
	}

	newChat := map[string]types.AttributeValue{
		"id":     &types.AttributeValueMemberS{Value: chatId},
		"active": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().UnixMilli())},
	}

	if len(chat.Messages) > 0 {
		newChat["messages"] = &types.AttributeValueMemberSS{Value: chat.Messages}
	}

	err = db.CreateChat(client, "chats", newChat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Chat opened!",
		"chat.id": chatId,
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

func CreateChatMessage(client *dynamodb.Client, w http.ResponseWriter, r *http.Request, chatId string) {

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

	var message db.Message
	err := json.NewDecoder(r.Body).Decode(&message)
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

	messageId := fmt.Sprintf("m_%s", id)

	newMessage := map[string]types.AttributeValue{
		"id":     &types.AttributeValueMemberS{Value: messageId},
		"sender": &types.AttributeValueMemberS{Value: message.Sender},
		"date":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", time.Now().UnixMilli())},
	}

	if message.Text != "" {
		newMessage["text"] = &types.AttributeValueMemberS{Value: message.Text}
	}

	if len(message.Media) > 0 {
		newMessage["media"] = &types.AttributeValueMemberSS{Value: message.Media}
	}

	err = db.CreateMessage(client, "messages", newMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := db.GetUserById(client, "chats", chatId)
	if err != nil {
		http.Error(w, `{"error": "Failed to get all users"}`, http.StatusInternalServerError)
		return
	}

	var chat db.Chat
	err = attributevalue.UnmarshalMap(resp, &chat)
	if err != nil {
		http.Error(w, `{"error": "Failed to decode chat"}`, http.StatusInternalServerError)
		return
	}

	chat.Messages = append(chat.Messages, messageId)
	err = db.UpdateChat(client, "chats", chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":    "Chat message sent!",
		"message.id": messageId,
		"chat.id":    chatId,
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

func GetChatById(client *dynamodb.Client, w http.ResponseWriter, r *http.Request, id string) {

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

	resp, err := db.GetChatById(client, "chats", id)
	if err != nil {
		http.Error(w, `{"error": "Failed to get chat"}`, http.StatusInternalServerError)
		return
	}

	var chat db.Chat
	err = attributevalue.UnmarshalMap(resp, &chat)
	if err != nil {
		http.Error(w, `{"error": "Failed to decode chat"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Got chat!",
		"chat":    chat,
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

func GetAllChats(client *dynamodb.Client, w http.ResponseWriter, r *http.Request) {

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

	resp, err := db.GetAllChats(client, "chats")
	if err != nil {
		http.Error(w, `{"error": "Failed to get all chats"}`, http.StatusInternalServerError)
		return
	}

	var chats []db.Chat
	err = attributevalue.UnmarshalListOfMaps(resp, &chats)
	if err != nil {
		http.Error(w, `{"error": "Failed to decode chats"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Got chat messages!",
		"chats":   chats,
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

func GetChatMessages(client *dynamodb.Client, w http.ResponseWriter, r *http.Request, chatId string) {

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

	resp, err := db.GetChatById(client, "chats", chatId)
	if err != nil {
		http.Error(w, `{"error": "Failed to get chat"}`, http.StatusInternalServerError)
		return
	}

	var chat db.Chat
	err = attributevalue.UnmarshalMap(resp, &chat)
	if err != nil {
		http.Error(w, `{"error": "Failed to decode chat"}`, http.StatusInternalServerError)
		return
	}

	var messages []db.Message
	for _, messageId := range chat.Messages {
		resp, err := db.GetMessageById(client, "messages", messageId)
		if err != nil {
			http.Error(w, `{"error": "Failed to get message"}`, http.StatusInternalServerError)
			return
		}

		var message db.Message
		err = attributevalue.UnmarshalMap(resp, &message)
		if err != nil {
			http.Error(w, `{"error": "Failed to decode message"}`, http.StatusInternalServerError)
			return
		}

		messages = append(messages, message)
	}

	response := map[string]interface{}{
		"message":  "Got chat messages!",
		"messages": messages,
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

func UpdateChat(client *dynamodb.Client, w http.ResponseWriter, r *http.Request) {

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

	var chat db.Chat
	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = db.UpdateChat(client, "chats", chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := `{"message": "Chat updated"}`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func DeleteChat(client *dynamodb.Client, w http.ResponseWriter, r *http.Request, id string) {

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

	err := db.DeleteChat(client, "chats", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := `{"message": "Chat deleted"}`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func DeleteChatMessage(client *dynamodb.Client, w http.ResponseWriter, r *http.Request, chatId, messageId string) {

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

	err := db.DeleteMessage(client, "messages", messageId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := db.GetChatById(client, "chats", chatId)
	if err != nil {
		http.Error(w, `{"error": "Failed to get chat"}`, http.StatusInternalServerError)
		return
	}

	var chat db.Chat
	err = attributevalue.UnmarshalMap(resp, &chat)
	if err != nil {
		http.Error(w, `{"error": "Failed to decode chat"}`, http.StatusInternalServerError)
		return
	}

	var newMessages []string
	for _, m := range chat.Messages {
		if m != messageId {
			newMessages = append(newMessages, m)
		}
	}

	chat.Messages = newMessages
	err = db.UpdateChat(client, "chats", chat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := `{"message": "Chat message deleted"}`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
