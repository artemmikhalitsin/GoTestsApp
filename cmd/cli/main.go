package main

import (
	"fmt"
	"log"
	"os"

	poker "github.com/artemmikhalitsin/GoTestsApp"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Please {Name} wins to record a win")

	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatalf("Problem initializing file system store: %v", err)
	}
	defer close()

	game := poker.NewCLI(store, os.Stdin, poker.BlindAlerterFunc(poker.StdOutAlerter))
	game.PlayPoker()
}
