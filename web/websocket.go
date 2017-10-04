package web

import (
	"github.com/Logiraptor/word-bot/core"
	"github.com/gorilla/websocket"
)

func SendMove(conn *websocket.Conn, move core.ScoredMove) error {
	return conn.WriteJSON(move)
}
