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

	game := poker.NewTexasHoldem(store, poker.BlindAlerterFunc(poker.StdOutAlerter))
	cli := poker.NewCLI(os.Stdin, os.Stdout, game)
	cli.PlayPoker()
}
