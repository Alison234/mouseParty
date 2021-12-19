package main

import (
	"net/http"
	"sync"

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

func serveWs(w http.ResponseWriter, r *http.Request, roomId string, mx *sync.Mutex) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("upgreder err %v", err)
		return
	}

	c := &connection{send: make(chan []byte, 256), ws: ws}
	s := subscription{conn: c, room: roomId, sessionId: sessionId}

	mx.Lock()
	sessionId = sessionId + 1
	defer mx.Unlock()

	h.register <- s
	go s.writePump()
	go s.readPump()
}
