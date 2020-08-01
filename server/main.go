package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/salmaanrizvi/crossword-party/bus"
)

/**
main.go
*/
const port int = 8000

func main() {
	router := gin.Default()

	hub := bus.NewHub()
	go hub.Run()

	router.GET("/", func(c *gin.Context) {
		c.String(200, "We got Gin")
	})

	router.GET("/ws", func(c *gin.Context) {
		wshandler(hub, c.Writer, c.Request)
	})

	router.RunTLS(
		fmt.Sprintf("localhost:%d", port),
		"./localhost/localhost.crt",
		"./localhost/localhost.key",
	)
}

func checkOrigin(r *http.Request) bool {
	return true
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func wshandler(hub *bus.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	client := bus.NewClient(hub, conn)

	go client.WritePump()
	go client.ReadPump()
}
