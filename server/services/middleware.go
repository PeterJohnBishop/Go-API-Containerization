package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *ResponseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriterWrapper) Write(data []byte) (int, error) {
	rw.body.Write(data)
	return rw.ResponseWriter.Write(data)
}

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &ResponseWriterWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           new(bytes.Buffer),
		}

		next(wrapper, r)

		duration := time.Since(start)
		logMsg := fmt.Sprintf("[%s] %s from %s - Status: %d, Duration: %s",
			r.Method, r.URL.Path, r.RemoteAddr, wrapper.statusCode, duration)

		if wrapper.statusCode >= 400 {
			logMsg += fmt.Sprintf(" - Response: %s", wrapper.body.String())
		}

		log.Println(logMsg)
	}
}

// authenticate JWT token
func VerifyJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix((r.URL.Path), "/register") || strings.HasPrefix((r.URL.Path), "/login") {
			next.ServeHTTP(w, r)
			return
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Authentication Header is missing!"}`, http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		userClaims := ParseAccessToken(token)
		if userClaims == nil {
			http.Error(w, `{"error": "Failed to verify token!"}`, http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// authenticate Refresh Token
type VerifyRefreshRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type ContextKey string

const userIDKey ContextKey = "userID"

func VerifyRefreshToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req VerifyRefreshRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		token := req.Token
		id := req.ID

		claims := ParseRefreshToken(token)
		if claims == nil {
			http.Error(w, `{"error": "Failed to verify token!"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, id)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
