package poker_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	poker "github.com/artemmikhalitsin/GoTestsApp"
	"github.com/gorilla/websocket"
)

var dummyGame = &GameSpy{}

func TestGETPlayers(t *testing.T) {
	store := poker.StubPlayerStore{
		Scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		WinCalls: []string{},
	}

	server, _ := poker.NewPlayerServer(&store, dummyGame)

	tests := []struct {
		name               string
		player             string
		expectedHTTPStatus int
		expectedScore      string
	}{
		{
			name:               "Return's Pepper's score",
			player:             "Pepper",
			expectedHTTPStatus: http.StatusOK,
			expectedScore:      "20",
		},
		{
			name:               "Returns Floyd's score",
			player:             "Floyd",
			expectedHTTPStatus: http.StatusOK,
			expectedScore:      "10",
		},
		{
			name:               "Returns 404 on missing players",
			player:             "Mario",
			expectedHTTPStatus: http.StatusNotFound,
			expectedScore:      "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := newGetScoreRequest(tt.player)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertResponseBody(t, response.Body.String(), tt.expectedScore)
			assertStatusCode(t, response.Code, tt.expectedHTTPStatus)
		})
	}
}

func TestStoreWins(t *testing.T) {
	store := poker.StubPlayerStore{
		Scores:   map[string]int{},
		WinCalls: []string{},
	}
	server, _ := poker.NewPlayerServer(&store, dummyGame)

	t.Run("it returns accepted on POST", func(t *testing.T) {
		name := "Pepper"
		response := httptest.NewRecorder()
		request := newPostWinRequest(name)

		server.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusAccepted)

		if len(store.WinCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.WinCalls), 1)
		}

		assertPlayerWin(t, store.WinCalls[0], name)
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := poker.League{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store := &poker.StubPlayerStore{League: wantedLeague}
		server := makePlayerServer(t, store, dummyGame)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)

		assertLeague(t, got, wantedLeague)
		assertStatusCode(t, response.Code, http.StatusOK)

		contentType := response.Result().Header.Get("content-type")
		assertContentType(t, contentType, jsonContentType)
	})
}

func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		store := &poker.StubPlayerStore{}
		server := makePlayerServer(t, store, dummyGame)

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusOK)
	})

	t.Run("when we get a message over a websocket, it's the winner", func(t *testing.T) {
		winner := "Wario"
		wantedBlindAlert := "Blind is 100"
		game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
		server := httptest.NewServer(makePlayerServer(t, dummyPlayerStore, game))
		ws := dialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")
		defer server.Close()
		defer ws.Close()

		//Start with 3 players
		writeWSMessage(t, ws, "3")
		//Send winner
		writeWSMessage(t, ws, winner)

		assertStartedWith(t, game, 3)
		assertFinishedWith(t, game, winner)

		within(t, 1*time.Second, func() { assertWebsocketGotMessage(t, ws, wantedBlindAlert) })
	})
}

func newGetScoreRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return request
}

func newPostWinRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return request
}

func newLeagueRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return request
}

func newGameRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return request
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("Got %q, want %q", got, want)
	}
}

func assertStatusCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}

func assertLeague(t *testing.T, got, wantedLeague poker.League) {
	t.Helper()
	if !reflect.DeepEqual(got, wantedLeague) {
		t.Errorf("Got %v, wanted %v", got, wantedLeague)
	}
}

const jsonContentType = "application/json"

func assertContentType(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response did not have content-type of %v, got %v", want, got)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one: %v", err)
	}
}

func assertPlayerWin(t *testing.T, got, want string) {
	if got != want {
		t.Fatalf("Expected to have winner %q, but got %q", want, got)
	}
}

func assertWebsocketGotMessage(t *testing.T, ws *websocket.Conn, want string) {
	_, got, _ := ws.ReadMessage()
	if string(got) != want {
		t.Errorf("Expected blind alert %q, but got %q", want, string(got))
	}
}

func getLeagueFromResponse(t *testing.T, body io.Reader) (league poker.League) {
	t.Helper()
	league, _ = poker.NewLeague(body)

	return
}

func makePlayerServer(t *testing.T, store poker.PlayerStore, game poker.Game) *poker.PlayerServer {
	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		t.Errorf("Error creating player server: %v", err)
	}
	return server
}

func dialWS(t *testing.T, wsURL string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Unable to open the websocket at %s: %v", wsURL, err)
	}

	return ws
}

func writeWSMessage(t *testing.T, ws *websocket.Conn, message string) {
	t.Helper()
	if err := ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("Unable to write message %s over the ws connection: %v", message, err)
	}
}

func within(t *testing.T, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("Function timed out")
	case <-done:
	}
}

func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)

	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}

	return false
}
