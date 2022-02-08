package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/kananb/chess"
)

type gameState struct {
	FEN        string
	SideToMove string
	History    []string
	WhiteClock int
	BlackClock int
	Loser      string
	moves      []chess.Move
	gameOver   bool

	latest  string
	current bool
}

func (g *gameState) String() string {
	if !g.current {
		s, _ := json.Marshal(*g)
		g.latest = string(s)
		g.current = true
	}

	return g.latest
}

type turnTimer struct {
	Color      chess.SideColor
	start, end int64
}

func (t *turnTimer) update(state *gameState, color chess.SideColor) bool {
	t.end = time.Now().UnixMilli()
	elapsed := int(t.end-t.start) / 100 // elapsed time in deciseconds
	t.start = t.end

	if t.Color == chess.White {
		state.WhiteClock -= elapsed
		if state.WhiteClock <= 0 {
			state.WhiteClock = 0
			state.Loser = chess.White.String()
			state.gameOver = true
			state.current = false
			return false
		}
	} else {
		state.BlackClock -= elapsed
		if state.BlackClock <= 0 {
			state.BlackClock = 0
			state.Loser = chess.White.String()
			state.gameOver = true
			state.current = false
			return false
		}
	}

	t.Color = color
	return true
}

func manageGame(room *gameRoom) {
	room.Board = chess.StartingPosition()
	state := gameState{WhiteClock: 69, BlackClock: 69}
	colors := [2]chess.SideColor{}

	rand.Seed(time.Now().UTC().UnixNano())
	colors[0] = chess.SideColor(rand.Intn(2) + 1)
	colors[1] = ^colors[0] & 3
	room.Users[0].send("STATE", fmt.Sprintf(`{"Side":%q}`, colors[0].String()))
	room.Users[1].send("STATE", fmt.Sprintf(`{"Side":%q}`, colors[1].String()))

	turn := turnTimer{room.Board.SideToMove, time.Now().UnixMilli(), 0}
gameLoop:
	for {
		if !state.current {
			state.moves = room.Board.Moves()
			if room.Board.InCheckmate() {
				state.Loser = room.Board.SideToMove.String()
				state.gameOver = true
			} else if room.Board.InStalemate() {
				state.Loser = "-"
				state.gameOver = true
			}

			turn.update(&state, room.Board.SideToMove)
			state.FEN = room.Board.String()
			state.SideToMove = room.Board.SideToMove.String()
			state.History = room.Board.History()

			go room.Broadcast("STATE", state.String())
		}

		var msg *message
		var from int
		for msg == nil {
			select {
			case msg = <-room.Users[0].ch:
				from = 0
			case msg = <-room.Users[1].ch:
				from = 1
			case <-time.After(time.Second):
				if ok := turn.update(&state, room.Board.SideToMove); !ok {
					continue gameLoop
				}
			}
		}
		if room.IsEmpty() {
			break
		}

		if msg.Cmd == "QUIT" {
			if !state.gameOver {
				state.gameOver = true
				state.Loser = colors[from].String()
				state.current = false
			}
		} else if msg.Cmd == "UPDATE" || state.gameOver {
			room.Users[from].send("STATE", fmt.Sprintf(`{"Side":%q}`, colors[from].String()))
			room.Users[from].send("STATE", state.String())
		} else if colors[from] != room.Board.SideToMove {
			room.Users[from].send("ERROR", "not your turn")
		} else if msg.Cmd == "MOVE" {
			move, err := chess.NewMove(msg.Args[0], room.Board)
			if err != nil {
				room.Users[from].send("ERROR", err.Error())
				continue
			}
			for _, legal := range state.moves {
				if legal.Matches(move) {
					if actual := room.Board.MakeMove(move); !actual.IsValid() {
						room.Users[from].send("ERROR", "move is invalid")
						continue gameLoop
					}

					break
				}
			}
			state.current = false
		}
	}
}
