package main

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

type hub struct {
	rooms      map[string]map[*connection]bool
	broadcast  chan message
	register   chan subscription
	unregister chan subscription
}

func (h *hub) run() {
	for {
		select {
		case s := <-h.register:
			fmt.Println(s.sessionId)
			connections := h.rooms[s.room]
			logrus.Info("New subscriber")

			registerMsg := RegisterMSG{
				Method:    "register",
				SessionId: s.sessionId,
				RoomId:    0,
			}

			data, err := json.Marshal(registerMsg)
			if err != nil {
				logrus.Info("failed send register msg to subscriber")
			}

			s.conn.send <- data

			if connections == nil {
				connections = make(map[*connection]bool)
				h.rooms[s.room] = connections
			}
			h.rooms[s.room][s.conn] = true

		case s := <-h.unregister:
			connections := h.rooms[s.room]
			logrus.Info("Out subscriber")

			outSubscribeMsg := RegisterMSG{
				Method:    "leave",
				SessionId: s.sessionId,
				RoomId:    0,
			}

			data, err := json.Marshal(outSubscribeMsg)
			if err != nil {
				logrus.Info("failed send register msg to subscriber")
			}

			s.conn.send <- data

			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, s.room)
					}
				}
			}
		case m := <-h.broadcast:

			connections := h.rooms[m.room]
			for c := range connections {
				select {
				case c.send <- m.data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, m.room)
					}
				}
			}
		}
	}
}
