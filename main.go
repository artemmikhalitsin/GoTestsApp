package main

import (
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(PlayerServer)
	err := http.ListenAndServe(":5000", handler)
	if err != nil {
		log.Fatalf("Could not listen on port 5000 %v", err)
	}
}
