package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const PlayerPrompt = "Please enter the number of players: "

type CLI struct {
	store   PlayerStore
	in      *bufio.Scanner
	out     io.Writer
	alerter BlindAlerter
}

func NewCLI(store PlayerStore, in io.Reader, out io.Writer, alerter BlindAlerter) *CLI {
	return &CLI{
		store:   store,
		in:      bufio.NewScanner(in),
		out:     out,
		alerter: alerter,
	}
}

func (cli *CLI) PlayPoker() {
	numPlayers, _ := strconv.Atoi(cli.readLine())
	cli.scheduleBlindAlerts(numPlayers)
	fmt.Fprintf(cli.out, PlayerPrompt)
	winner := extractWinner(cli.readLine())
	cli.store.RecordWin(winner)
}

func extractWinner(input string) string {
	return strings.Replace(input, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

func (cli *CLI) scheduleBlindAlerts(numPlayers int) {
	blindIncrement := time.Duration(5+numPlayers) * time.Minute
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += blindIncrement
	}
}
