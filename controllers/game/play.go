package game

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/kananb/chess"
)

type gameState struct {
	FEN        string
	SideToMove string
	History    []string
	Loser      string
}

func awaitMessage(room *gameRoom) (msg *message, from int) {
	for msg == nil {
		from = ^from & 1
		if room.IsEmpty() {
			return nil, 0
		} else if room.players[from] == nil {
			continue
		}

		select {
		case msg = <-room.players[from].ch:
		default:
		}
	}

	return
}

func manageGame(room *gameRoom) {
	room.Board = chess.StartingPosition()
	var moves []chess.Move

	var latestState []byte
	current := false
	loser := chess.SideColor(0)

	colors := [2]chess.SideColor{}
	colors[0] = chess.SideColor(rand.Intn(2) + 1)
	colors[1] = ^colors[0] & 3

	tellColors := func(i int) {
		if room.players[i] != nil {
			communicator{room.players[i].Conn}.send("STATE", fmt.Sprintf(`{"Side":%q}`, colors[i]))
		}
	}
	tellColors(0)
	tellColors(1)

gameLoop:
	for {
		if !current {
			moves = room.Board.Moves()
			if room.Board.InCheckmate() {
				loser = room.Board.SideToMove
			}

			state, err := json.Marshal(gameState{
				room.Board.String(),
				room.Board.SideToMove.String(),
				room.Board.History(),
				loser.String(),
			})
			if err != nil {
				fmt.Println(err)
				break
			}

			go room.Broadcast("STATE", string(state))
			latestState = state
			current = true
		}

		msg, from := awaitMessage(room)
		if msg == nil {
			break
		}

		if (loser.IsValid() || msg.Cmd == "UPDATE") && msg.Cmd != "QUIT" {
			tellColors(from)
			communicator{room.players[from].Conn}.send("STATE", string(latestState))
		} else if msg.Cmd == "QUIT" {
			loser = colors[from]
			current = false
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
