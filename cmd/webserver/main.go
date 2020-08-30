package main

import (
	"log"
	"net/http"

	poker "github.com/artemmikhalitsin/GoTestsApp"
)

const dbFileName = "game.db.json"

func main() {

	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatalf("problem creating filesystem store: %v", err)
	}
	defer close()

	scoreServer := poker.NewPlayerServer(store)

	err = http.ListenAndServe(":5000", scoreServer)
	if err != nil {
		log.Fatalf("Could not listen on port 5000 %v", err)
	}
}
