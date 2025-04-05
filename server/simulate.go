package server

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"upgraded-telegram/main.go/server/services/db"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func generateRandomUserPayload() []byte {
	u := db.User{
		ID:       fmt.Sprintf("u_%s", gofakeit.UUID()),
		Name:     fmt.Sprintf("%s %s", gofakeit.FirstName(), gofakeit.Username()),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, false, false, 8),
	}
	jsonPayload, _ := json.Marshal(u)
	return jsonPayload
}

func SimulateRequests() {
	gofakeit.Seed(time.Now().UnixNano())

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwt := os.Getenv("SIMULATOR_JWT")

	rate := vegeta.Rate{Freq: 100, Per: time.Second} // 100 RPS
	duration := 10 * time.Second
	attacker := vegeta.NewAttacker()

	targeter := func(t *vegeta.Target) error {
		t.Method = "POST"
		t.URL = "http://0.0.0.0:8080/users/new"
		t.Header = map[string][]string{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer " + jwt},
		}
		t.Body = generateRandomUserPayload()
		return nil
	}

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Randomized User Test") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Println("Requests:", metrics.Requests)
	fmt.Printf("Success Rate: %.2f%%\n", metrics.Success*100)
	fmt.Println("Avg Latency:", metrics.Latencies.Mean)
	fmt.Println("99th Percentile:", metrics.Latencies.P99)
	fmt.Println("Errors:", metrics.Errors)
}
