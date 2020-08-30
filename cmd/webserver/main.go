package main

import (
	"log"
	"net/http"
	"os"

	poker "github.com/artemmikhalitsin/GoTestsApp"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("Problem opening %s: %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating filesystem store: %v", err)
	}

	scoreServer := poker.NewPlayerServer(store)

	err = http.ListenAndServe(":5000", scoreServer)
	if err != nil {
		log.Fatalf("Could not listen on port 5000 %v", err)
	}
}
