package main

type message struct {
	data []byte
	room string
	Id   int
}

type RegisterMSG struct {
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Method    string `json:"method"`
	SessionId int    `json:"sessionId"`
	RoomId    int    `json:"roomId"`
}
