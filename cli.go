package poker

import (
	"bufio"
	"io"
	"strings"
	"time"
)

type CLI struct {
	store   PlayerStore
	in      *bufio.Scanner
	alerter BlindAlerter
}

func NewCLI(store PlayerStore, in io.Reader, alerter BlindAlerter) *CLI {
	return &CLI{
		store:   store,
		in:      bufio.NewScanner(in),
		alerter: alerter,
	}
}

func (cli *CLI) PlayPoker() {
	cli.scheduleBlindAlerts()
	winner := extractWinner(cli.readLine())
	cli.store.RecordWin(winner)
}

func extractWinner(input string) string {
	return strings.Replace(input, " wins", "", 1)
}

type BlindAlerter interface {
	ScheduleAlertAt(time time.Duration, value int)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

func (cli *CLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += 10 * time.Minute
	}
}
