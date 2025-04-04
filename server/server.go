package server

import (
	"fmt"
	"log"
	"net/http"
	"upgraded-telegram/main.go/server/handlers"
	"upgraded-telegram/main.go/server/services"
	"upgraded-telegram/main.go/server/services/ai"
	"upgraded-telegram/main.go/server/services/db"
	"upgraded-telegram/main.go/server/services/fileIO"
	"upgraded-telegram/main.go/server/services/mapping"

	"googlemaps.github.io/maps"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/openai/openai-go"
)

func StartServer() {

	mux := http.NewServeMux()

	// implement rate limiting
	rateLimiter := services.NewRateLimiter(5, 10)
	handler := rateLimiter.RateLimitMiddleware(mux)

	// connect with AWS
	cfg := services.StartAws()

	// connect with DynamoDB
	dynamoClient := db.ConnectDB(cfg)

	// connect with S3
	s3Client := fileIO.ConnectS3(cfg)

	// connect with OpenAI
	aiClient := ai.Open()

	// connect with Google Maps
	mapClient := mapping.FindMaps()

	services.InitAuth()
	addUserRoutes(dynamoClient, mux)
	addChatMessageRoutes(dynamoClient, mux)
	addFileIORoutes(s3Client, mux)
	addAIRoutes(aiClient, mux)
	addEventRoutes(dynamoClient, mux)
	addMapRoutes(mapClient, mux)

	fmt.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatalf("unable to load dynamoDB tables, %v", err)
	}
}

func addUserRoutes(client *dynamodb.Client, mux *http.ServeMux) {
	mux.HandleFunc("/users/new", services.LoggerMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateUser(client, w, r)
	}))
	mux.HandleFunc("/users/login", services.LoggerMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.AuthUser(client, w, r)
	}))
	mux.HandleFunc("/users/all", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllUsers(client, w, r)
	})))
	mux.HandleFunc("/users/id/{id}", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		handlers.GetUserByID(client, w, r, id)
	})))
	mux.HandleFunc("/users/update", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateUser(client, w, r)
	})))
	mux.HandleFunc("/users/delete/{id}", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		handlers.DeleteUser(client, w, r, id)
	})))
}

func addChatMessageRoutes(client *dynamodb.Client, mux *http.ServeMux) {
	mux.HandleFunc("/chats/new", services.LoggerMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateChat(client, w, r)
	}))
	mux.HandleFunc("/chats/chat/{id}/messages/new", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		handlers.CreateChatMessage(client, w, r, id)
	})))
	mux.HandleFunc("/chats/all", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllChats(client, w, r)
	})))
	mux.HandleFunc("/chats/chat/{id}", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		handlers.GetChatById(client, w, r, id)
	})))
	mux.HandleFunc("/chats/chat/{chatId}/messages", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("chatId")
		handlers.GetChatMessages(client, w, r, id)
	})))
	mux.HandleFunc("/chats/chat/update", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateChat(client, w, r)
	})))
	mux.HandleFunc("/chats/chat/{id}/delete", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		handlers.DeleteChat(client, w, r, id)
	})))
	mux.HandleFunc("/chats/chat/{chatId}/messages/message/{messageId}/delete", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		chatId := r.PathValue("chatId")
		messageId := r.PathValue("messageId")
		handlers.DeleteChatMessage(client, w, r, chatId, messageId)
	})))
}

func addFileIORoutes(client *s3.Client, mux *http.ServeMux) {
	mux.HandleFunc("/upload", services.LoggerMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleFileUpload(client, w, r)
	}))
	mux.HandleFunc("/download", services.LoggerMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleFileDownload(client, w, r)
	}))
}

func addAIRoutes(client *openai.Client, mux *http.ServeMux) {
	// ready for routes!
}

func addEventRoutes(client *dynamodb.Client, mux *http.ServeMux) {
	mux.HandleFunc("/events/new", services.LoggerMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateEventHandler(client, w, r)
	}))
	mux.HandleFunc("/events/event/{id}", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		handlers.GetEventById(client, w, r, id)
	})))
	mux.HandleFunc("/events/all", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAllEvents(client, w, r)
	})))
	mux.HandleFunc("/events/event/update", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateChat(client, w, r)
	})))
	mux.HandleFunc("/events/event/{id}/delete", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		handlers.DeleteEvent(client, w, r, id)
	})))
}

func addMapRoutes(client *maps.Client, mux *http.ServeMux) {
	mux.HandleFunc("/maps/from/{origin}/to/{destination}", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		a := r.PathValue("origin")
		b := r.PathValue("destination")
		handlers.GetDirections(client, w, r, a, b)
	})))
	mux.HandleFunc("/maps/reversegeocode/lat/{lat}/long/{long}", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		long := r.PathValue("long")
		lat := r.PathValue("lat")
		handlers.ReverseGeocode(client, w, r, lat, long)
	})))
	mux.HandleFunc("/maps/geocode/{address}", services.LoggerMiddleware(services.VerifyJWT(func(w http.ResponseWriter, r *http.Request) {
		address := r.PathValue("address")
		handlers.Geocode(client, w, r, address)
	})))
}
