package poker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const PlayerPrompt = "Please enter the number of players: "
const BadInputErrMessage = "Can't understand input"
const BadWinnerInputErrMessage = "Winner message doesn't look right"

// CLI is the command-line program to run a poker game
type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game Game
}

// NewCLI creates a new CLI object given input source, output destination and a game to control
func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

// PlayPoker starts interaction with the user to run the poker game
func (cli *CLI) PlayPoker() {
	fmt.Fprintf(cli.out, PlayerPrompt)

	numPlayersInput := cli.readLine()
	numPlayers, err := strconv.Atoi(numPlayersInput)

	if err != nil {
		fmt.Fprintf(cli.out, BadInputErrMessage)
		return
	}

	cli.game.Start(numPlayers, cli.out)

	winnierInput := cli.readLine()
	winner, err := extractWinner(winnierInput)

	if err != nil {
		fmt.Fprintf(cli.out, BadWinnerInputErrMessage)
		return
	}

	cli.game.Finish(winner)
}

func extractWinner(input string) (string, error) {
	if !strings.Contains(input, " wins") {
		return "", errors.New(BadWinnerInputErrMessage)
	}
	return strings.Replace(input, " wins", "", 1), nil
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
