package game

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kananb/blind-chess/data"
	"github.com/kananb/chess"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(*http.Request) bool { return true },
}

var manager data.ChessManager

func init() {
	if _, present := os.LookupEnv("REDISHOST"); present {
		manager = data.NewChessManager(data.NewRedisStateManager())
	} else {
		manager = data.NewChessManager(data.NewMemoryStateManager())
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
	var code, id string
	var in chan string
	defer func() {
		manager.Leave(code, in)
		conn.Close()
	}()

	msgs, msg, ok := make(chan *message), (*message)(nil), false
	go func() {
		for {
			msg, err := comm.receive()
			if err != nil {
				fmt.Println(err)
				close(msgs)
				return
			}

			msgs <- msg
		}
	}()

awaitGame:
	// wait for connection to join or create a game room
	code = ""
	for code == "" {
		if msg, ok = <-msgs; !ok {
			return
		}

		id = ""
		if msg.Cmd == "JOIN" {
			if msg.Args[0] == "" {
				comm.send("DENY", "no room code provided")
				continue
			}

			if len(msg.Args) >= 2 {
				id = msg.Args[1]
			}
			if id, in, err = manager.Join(msg.Args[0], id); err != nil {
				comm.send("DENY", err.Error())
				continue
			}
			code = msg.Args[0]
		} else if msg.Cmd == "CREATE" {
			config := data.NewGameConfig(msg.Args[0])
			if config == nil {
				comm.send("DENY", "invalid game configuration")
				continue
			}
			code, id, in = manager.Create(*config)
			if in == nil {
				comm.send("DENY", "something went wrong")
			}
		}
	}
	comm.send("IN", code, id)
	state, side := manager.Get(code), chess.White
	if state.Players[1].ID == id {
		side = chess.Black
	}
	var board *chess.Board = nil
	if state.FEN != "" {
		board, _ = chess.NewBoard(state.FEN)
	}

	updateTime := func() {
		diff := int(time.Now().UnixMilli()-state.TimeOfLastMove) / 100
		if board.SideToMove == chess.White {
			state.Players[0].Clock -= diff
			if state.Players[0].Clock <= 0 {
				state.Players[0].Clock = 0
				state.Result = "0-1"
			}
		} else {
			state.Players[1].Clock -= diff
			if state.Players[1].Clock <= 0 {
				state.Players[1].Clock = 0
				state.Result = "1-0"
			}
		}
		state.TimeOfLastMove = time.Now().UnixMilli()
	}
	updateResult := func(won bool) {
		if (side == chess.White && !won) || (side == chess.Black && won) {
			state.Result = "0-1"
		} else {
			state.Result = "1-0"
		}
	}
	notify := func() {
		manager.Notify(state, code, in)
		comm.send("STATE", state.Marshal(id))
	}

gameLoop:
	for {
		select {
		case msg, ok = <-msgs:
			if !ok {
				return
			}
		case <-in:
			state = manager.Get(code)
			board, _ = chess.NewBoard(state.FEN)
			comm.send("STATE", state.Marshal(id))
			continue
		}

		if msg.Cmd == "QUIT" {
			if board != nil && state.Result == "" {
				updateTime()
				updateResult(false)
				manager.Notify(state, code, in)
			}
			break
		} else if msg.Cmd == "UPDATE" {
			if board != nil {
				updateTime()
				manager.Notify(state, code, in)
			}
			comm.send("STATE", state.Marshal(id))
			continue
		} else if msg.Cmd == "JOIN" || msg.Cmd == "CREATE" {
			comm.send("IN", code, id)
			continue
		} else if board == nil {
			comm.send("ERROR", "game hasn't started yet")
			continue
		} else if msg.Cmd == "RESIGN" {
			if board != nil && state.Result == "" {
				updateTime()
				updateResult(false)
				notify()
			} else {
				comm.send("STATE", state.Marshal(id))
			}
		} else if msg.Cmd == "MOVE" {
			if side != board.SideToMove {
				comm.send("ERROR", "not your turn")
				continue
			}

			move, err := chess.NewMove(msg.Args[0], board)
			if err != nil {
				comm.send("ERROR", err.Error())
				continue
			}

			updateTime()
			var actual chess.Move
			for _, legal := range board.Moves() {
				if !legal.Matches(move) {
					continue
				} else if actual = board.MakeMove(move); !actual.IsValid() {
					comm.send("ERROR", "move is invalid")
					continue gameLoop
				}
				break
			}

			if board.SideToMove == chess.White {
				state.Players[1].Clock += state.Increment
			} else {
				state.Players[0].Clock += state.Increment
			}

			if board.InCheckmate() {
				actual.Check = chess.Checkmate
				updateResult(true)
			} else if board.InCheck(board.SideToMove) {
				actual.Check = chess.Check
			} else if board.InStalemate() {
				state.Result = "1/2-1/2"
			}
			state.FEN = board.String()
			state.SideToMove = board.SideToMove
			state.History = append(state.History, actual.String())
			state.TimeOfLastMove = time.Now().UnixMilli()

			notify()
		}
	}

	manager.Leave(code, in)
	goto awaitGame
}

func Route(router *gin.RouterGroup) {
	router.GET("", handleWebsocket)
}
