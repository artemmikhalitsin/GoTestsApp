package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	store PlayerStore
	in    *bufio.Scanner
}

func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{
		store: store,
		in:    bufio.NewScanner(in),
	}
}

func (cli *CLI) PlayPoker() {
	cli.in.Scan()
	winner := extractWinner(cli.in.Text())
	cli.store.RecordWin(winner)
}

func extractWinner(input string) string {
	return strings.Replace(input, " wins", "", 1)
}
