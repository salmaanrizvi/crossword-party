package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/salmaanrizvi/crossword-party/server/bus"
	"github.com/salmaanrizvi/crossword-party/server/config"
)

func main() {
	conf := config.Get()
	defer config.Logger().Sync()

	if conf.Env == config.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(
		config.GetGinLoggerMiddleware(config.Logger()),
		// config.GetRecoveryLoggerMiddleware(config.Logger(), true),
	)

	hub := bus.NewHub()
	go hub.Run()
	go hub.Stats(conf.LogStatsInterval)

	router.GET("/", func(c *gin.Context) {
		c.String(200, fmt.Sprintf("Crossword Party %s", conf.AppVersion))
	})

	router.GET("/ws", func(c *gin.Context) {
		wshandler(hub, c.Writer, c.Request)
	})

	if conf.RunTLS() {
		router.RunTLS(
			fmt.Sprintf(":%d", conf.Port),
			conf.CertFile,
			conf.KeyFile,
		)
	} else {
		router.Run(fmt.Sprintf(":%d", conf.Port))
	}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	return origin == "https://www.nytimes.com"
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func wshandler(hub *bus.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		config.Logger().Errorf("Failed to set websocket upgrade %s", err.Error())
		return
	}

	client := bus.NewClient(hub, conn)

	go client.WritePump()
	go client.ReadPump()
}
