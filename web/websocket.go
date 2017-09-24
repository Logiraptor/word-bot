package web

import (
	"word-bot/ai"

	"github.com/gorilla/websocket"
)

func SendMove(conn *websocket.Conn, move ai.ScoredMove) error {
	return conn.WriteJSON(move)
}
