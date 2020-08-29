package main

import (
	"log"
	"net/http"
)

func main() {
	scoreServer := &PlayerServer{NewInMemoryPlayerStore()}
	err := http.ListenAndServe(":5000", scoreServer)
	if err != nil {
		log.Fatalf("Could not listen on port 5000 %v", err)
	}
}
