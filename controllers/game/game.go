package game

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

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

type player struct {
	Conn *websocket.Conn
	Side chess.SideColor
}

func (p *player) Exists() bool {
	return p.Conn != nil
}

type gameRoom struct {
	p1, p2 player
	board  *chess.Board
}

var rooms = map[string]*gameRoom{}
var roomMutex = sync.Mutex{}

func genRoomCode() string {
	return fmt.Sprintf("%04d", len(rooms))
}
func joinRoom(code string, conn *websocket.Conn) (string, error) {
	roomMutex.Lock()
	defer roomMutex.Unlock()

	if room, ok := rooms[code]; ok {
		if !room.p1.Exists() || !room.p2.Exists() {
			if !room.p1.Exists() {
				room.p1.Conn = conn
			} else {
				room.p2.Conn = conn
			}
		} else {
			return "", fmt.Errorf("room is already full")
		}
	} else {
		return "", fmt.Errorf("room doesn't exist")
	}

	return code, nil
}
func createRoom(conn *websocket.Conn) (code string) {
	roomMutex.Lock()
	defer roomMutex.Unlock()

	code = genRoomCode()
	rooms[code] = &gameRoom{
		player{
			Conn: conn,
		},
		player{},
		chess.StartingPosition(),
	}

	return
}
func leaveRoom(code string, conn *websocket.Conn) {
	roomMutex.Lock()
	defer roomMutex.Unlock()

	room, ok := rooms[code]
	if !ok {
		return
	}

	if room.p1.Conn == conn {
		room.p1.Conn = nil
	} else if room.p2.Conn == conn {
		room.p2.Conn = nil
	}
}

type message struct {
	Cmd  string
	Args []string
}

func newMessage(data []byte) *message {
	msg := new(message)
	parts := strings.Split(string(data), "_")

	msg.Cmd = parts[0]
	if len(parts) > 1 {
		msg.Args = parts[1:]
	}

	return msg
}

func handleWebsocket(c *gin.Context) {
	w, r := c.Writer, c.Request
	conn, err := upgrader.Upgrade(w, r, nil)
	var code string
	if err != nil {
		log.Printf("Failed to upgrade connection: %+v", err)
		return
	}

	defer conn.Close()
	defer func() {
		leaveRoom(code, conn)
	}()

	send := func(txt string) {
		conn.WriteMessage(websocket.TextMessage, []byte(txt))
	}

awaitGame:
	for code == "" {
		_, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		msg := newMessage(data)
		if msg.Cmd == "JOIN" {
			if len(msg.Args) > 0 {
				code, err = joinRoom(msg.Args[0], conn)
				if err != nil {
					send(fmt.Sprintf("DENY_%v", err))
				}
			} else {
				send("DENY_Must provide a room code")
			}
		} else if msg.Cmd == "CREATE" {
			code = createRoom(conn)
		}
	}

	send("CODE_" + code)

	for {
		t, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		conn.WriteMessage(t, data)
	}

	goto awaitGame
}

func Route(router *gin.RouterGroup) {
	router.GET("", handleWebsocket)
}
