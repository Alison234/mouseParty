package main

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var h = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[*connection]bool),
}

func main() {
	go h.run()
	mx := sync.Mutex{}
	router := gin.New()
	router.LoadHTMLFiles("..//client/index.html")

	router.GET("/room/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.GET("/ws", func(c *gin.Context) {
		serveWs(c.Writer, c.Request, "0", &mx)
	})

	logrus.Info("Staring server")

	err := router.Run("localhost:4567")
	if err != nil {
		logrus.Errorf("route server err %v", err)
	}
}
