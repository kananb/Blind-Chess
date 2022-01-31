package game

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// remove this function assignment outside of testing
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebsocket(c *gin.Context) {
	w, r := c.Writer, c.Request
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		conn.WriteMessage(t, msg)
	}
}

func Route(router *gin.RouterGroup) {
	router.GET("", handleWebsocket)
}
