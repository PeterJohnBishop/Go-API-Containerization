package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"upgraded-telegram/main.go/server/services"
	"upgraded-telegram/main.go/server/services/fileIO"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func HandleFileUpload(client *s3.Client, w http.ResponseWriter, r *http.Request) {

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

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileURL, err := fileIO.UploadFile(client, header.Filename, file)
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("File uploaded successfully: %s", fileURL)))
}

func HandleFileDownload(client *s3.Client, w http.ResponseWriter, r *http.Request) {

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

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Missing filename parameter", http.StatusBadRequest)
		return
	}

	url, err := fileIO.DownloadFile(client, filename)
	if err != nil {
		http.Error(w, "Failed to generate download URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"downloadUrl": "%s"}`, url)))
}
