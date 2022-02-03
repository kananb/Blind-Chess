package game

import (
	"fmt"
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

var manager = roomManager{}

func handleWebsocket(c *gin.Context) {
	w, r := c.Writer, c.Request
	conn, err := upgrader.Upgrade(w, r, nil)
	comm := communicator{conn}
	var code string
	if err != nil {
		log.Printf("Failed to upgrade connection: %+v", err)
		return
	}

	defer conn.Close()
	defer func() {
		manager.RemoveConn(conn, code)
	}()

awaitGame:
	// wait for connection to join or create a game room
	code = ""
	for code == "" {
		msg, err := comm.receive()
		if err != nil {
			fmt.Println(err)
			return
		}

		if msg.Cmd == "JOIN" {
			if msg.Args[0] == "" {
				comm.send("DENY", "no room code provided")
				continue
			}
			if err = manager.AddConn(conn, msg.Args[0]); err != nil {
				comm.send("DENY", err.Error())
				continue
			}
			code = msg.Args[0]
		} else if msg.Cmd == "CREATE" {
			code = manager.CreateRoom()
			if err = manager.AddConn(conn, code); err != nil {
				manager.RemoveRoom(code)
				comm.send("DENY", err.Error())
				code = ""
			}
		}
	}

	// Wait for game room code acknowledgement
	for {
		comm.send("CODE", code)
		msg, err := comm.receive()
		if err != nil {
			fmt.Println(err)
			return
		}

		if msg.Cmd == "OK" {
			break
		}
	}

	// handle game room communications
	for {
		msg, err := comm.receive()
		if err != nil {
			fmt.Println(err)
			return
		}

		if msg.Cmd == "LEAVE" {
			break
		}
		comm.send(msg.Cmd, msg.Args...)
	}

	goto awaitGame
}

func Route(router *gin.RouterGroup) {
	router.GET("", handleWebsocket)
}
