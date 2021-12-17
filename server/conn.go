package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var sessionId = 0

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type connection struct {
	ws   *websocket.Conn
	send chan []byte
}

func serveWs(w http.ResponseWriter, r *http.Request, roomId string) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("upgreder err %v", err)
		return
	}

	c := &connection{send: make(chan []byte, 256), ws: ws}
	s := subscription{conn: c, room: roomId, sessionId: sessionId}
	sessionId = sessionId + 1
	h.register <- s
	go s.writePump()
	go s.readPump()
}
