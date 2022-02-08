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

	CheckOrigin: func(*http.Request) bool { return true },
}

var manager = newRoomManager()

func handleWebsocket(c *gin.Context) {
	w, r := c.Writer, c.Request
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %+v", err)
		return
	}
	comm := communicator{conn}
	var code string

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
			code = manager.CreateRoom(msg.Args[0], manageGame)
			if err = manager.AddConn(conn, code); err != nil {
				manager.RemoveRoom(code)
				comm.send("DENY", err.Error())
				code = ""
			}
		}
	}
	comm.send("CODE", code)

	room := manager.get(code)
	ch := room.GetChannel(conn)
	if room.Started {
		ch <- &message{Cmd: "UPDATE"}
	}

	for {
		msg, err := comm.receive()
		if err != nil {
			fmt.Println(err)
			return
		}

		if room.Started {
			ch <- msg
		} else if msg.Cmd != "QUIT" {
			comm.send("ERROR", "game hasn't started yet")
		}
		if msg.Cmd == "QUIT" {
			break
		} else if msg.Cmd == "JOIN" || msg.Cmd == "CREATE" {
			comm.send("CODE", code)
		}
	}

	manager.RemoveConn(conn, code)
	goto awaitGame
}

func Route(router *gin.RouterGroup) {
	router.GET("", handleWebsocket)
}
