package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	cb := NewCircuitBreaker(10, 8, 5*time.Second)
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	for i := 0; i < 10; i++ {
		go func() {
			for {
				start := time.Now()
				res, err := cb.Call(func() (*http.Response, error) {
					return client.Get("http://localhost:8080/process")
				})
				duration := time.Since(start)

				if err != nil {
					log.Printf("Error: %v\n", err)
				} else {
					body, err := io.ReadAll(res.Body)
					if err == nil {
						log.Printf("Response in %dms: %s\n", duration.Milliseconds(), body)
					} else {
						log.Printf("Error reading response body: %v\n", err)
					}
					res.Body.Close()
				}

				time.Sleep(time.Duration(30+rand.Intn(70)) * time.Millisecond)
			}
		}()
	}
}
