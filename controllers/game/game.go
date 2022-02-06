package game

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kananb/chess"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// remove this function assignment outside of testing
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var manager = newRoomManager()

func manageGame(room *gameRoom) {
	room.Board = chess.StartingPosition()
	var moves []chess.Move
	current := false

	type boardInfo struct {
		FEN        string
		SideToMove string
		History    []string
	}

	colors := [2]chess.SideColor{}
	colors[0] = chess.SideColor(rand.Intn(2) + 1)
	colors[1] = ^colors[0] & 3

gameLoop:
	for {
		if !current {
			state, err := json.Marshal(boardInfo{
				room.Board.String(),
				room.Board.SideToMove.String(),
				room.Board.History(),
			})
			if err != nil {
				fmt.Println(err)
				break
			}
			go room.Broadcast("STATE", string(state))

			moves = room.Board.Moves()
			current = true
		}

		var msg *message
		from := 0
		for msg == nil {
			from = ^from & 1
			if room.IsEmpty() {
				break gameLoop
			} else if room.players[from] == nil {
				continue
			}

			select {
			case msg = <-room.players[from].ch:
			default:
			}
		}

		if msg.Cmd == "QUIT" {
			go room.Broadcast("END")
		} else if colors[from] != room.Board.SideToMove {
			communicator{room.players[from].Conn}.send("ERROR", "not your turn")
		} else if msg.Cmd == "MOVE" {
			move, err := chess.NewMove(msg.Args[0], room.Board)
			if err != nil {
				communicator{room.players[from].Conn}.send("ERROR", err.Error())
				continue
			}
			for _, legal := range moves {
				if legal.Matches(move) {
					if actual := room.Board.MakeMove(move); !actual.IsValid() {
						communicator{room.players[from].Conn}.send("ERROR", "move is invalid")
						continue gameLoop
					}
					break
				}
			}
			current = false
		}
	}
}

func handleWebsocket(c *gin.Context) {
	w, r := c.Writer, c.Request
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %+v", err)
		return
	}
	comm := communicator{conn}
	p := &player{Conn: conn}
	var code string

	defer conn.Close()
	defer func() {
		manager.RemovePlayer(p, code)
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
			if err = manager.AddPlayer(p, msg.Args[0]); err != nil {
				comm.send("DENY", err.Error())
				continue
			}
			code = msg.Args[0]
		} else if msg.Cmd == "CREATE" {
			code = manager.CreateRoom(manageGame)
			if err = manager.AddPlayer(p, code); err != nil {
				manager.RemoveRoom(code)
				comm.send("DENY", err.Error())
				code = ""
			}
		}
	}

	// wait for game room code acknowledgement
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

	room := manager.get(code)
	for {
		msg, err := comm.receive()
		if err != nil {
			fmt.Println(err)
			return
		}

		if room.Started {
			p.ch <- msg
		} else if msg.Cmd != "QUIT" {
			comm.send("ERROR", "game hasn't started yet")
		}
		if msg.Cmd == "QUIT" {
			break
		}
	}

	manager.RemovePlayer(p, code)
	goto awaitGame
}

func Route(router *gin.RouterGroup) {
	router.GET("", handleWebsocket)
}
