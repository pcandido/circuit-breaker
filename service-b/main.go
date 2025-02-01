package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"
)

type Config struct {
	SuccessRate int `json:"successRate"`
	LatencyMin  int `json:"latencyMin"`
	LatencyMax  int `json:"latencyMax"`
}

var config atomic.Value

func main() {
	config.Store(&Config{
		SuccessRate: 99,
		LatencyMin:  0,
		LatencyMax:  200,
	})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/healthcheck", healthCheckHandler)
	http.HandleFunc("/process", processHandler)
	http.HandleFunc("/config", configHandler)

	log.Println("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	currentConfig := config.Load().(*Config)

	latency := rand.Intn(currentConfig.LatencyMax-currentConfig.LatencyMin+1) + currentConfig.LatencyMin
	time.Sleep(time.Duration(latency) * time.Millisecond)

	if rand.Intn(100) < currentConfig.SuccessRate {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		log.Printf("OK: %dms\n", latency)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		log.Printf("ERROR: %dms\n", latency)
	}
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newConfig Config
	err := json.NewDecoder(r.Body).Decode(&newConfig)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	config.Store(&newConfig)

	log.Printf("Received config: %+v\n", newConfig)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
