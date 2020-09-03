package poker

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type playerServerWS struct {
	*websocket.Conn
}

func newPlayerServerWS(w http.ResponseWriter, r *http.Request) *playerServerWS {
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Problem upgrading connection to websocket: %v", err)
	}

	return &playerServerWS{conn}
}

func (w *playerServerWS) WaitForMsg() string {
	_, msg, err := w.Conn.ReadMessage()

	if err != nil {
		log.Printf("Problem reading message from websocker: %v", err)
	}

	return string(msg)
}

func (w *playerServerWS) Write(msg []byte) (int, error) {
	err := w.Conn.WriteMessage(websocket.TextMessage, msg)

	if err != nil {
		log.Printf("Error writing message to websocket: %v", err)
		return 0, err
	}

	return len(msg), nil
}
