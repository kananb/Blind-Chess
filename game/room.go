package game

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kananb/chess"
)

type user struct {
	communicator
	ch chan *message
}
type gameConfig struct {
	Duration  int
	Increment int
	PlayAs    string
}
type gameReadyCallabck func(*gameRoom)
type gameRoom struct {
	Users   [2]user
	Board   *chess.Board
	Config  gameConfig
	OnReady gameReadyCallabck
	Started bool
}

func (r *gameRoom) Join(c *websocket.Conn) error {
	var i int
	if r.Users[0].conn == nil {
		i = 0
	} else if r.Users[1].conn == nil {
		i = 1
	} else {
		return fmt.Errorf("room is full")
	}

	r.Users[i].conn = c
	return nil
}
func (r *gameRoom) Leave(c *websocket.Conn) {
	var i int
	if r.Users[0].conn == c {
		i = 0
	} else if r.Users[1].conn == c {
		i = 1
	} else {
		return
	}

	r.Users[i].conn = nil
}
func (r *gameRoom) GetChannel(c *websocket.Conn) chan *message {
	var i int
	if r.Users[0].conn == c {
		i = 0
	} else if r.Users[1].conn == c {
		i = 1
	} else {
		return nil
	}

	return r.Users[i].ch
}
func (r *gameRoom) Broadcast(cmd string, args ...string) {
	for _, p := range r.Users {
		p.send(cmd, args...)
	}
}

func (r *gameRoom) IsEmpty() bool {
	return r.Users[0].conn == nil && r.Users[1].conn == nil
}
func (r *gameRoom) IsFull() bool {
	return r.Users[0].conn != nil && r.Users[1].conn != nil
}

type roomManager struct {
	rooms map[string]*gameRoom
	lock  sync.RWMutex
}

func newRoomManager() (manager *roomManager) {
	manager = new(roomManager)
	manager.rooms = map[string]*gameRoom{}

	return
}

func (m *roomManager) get(code string) *gameRoom {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.rooms[code]
}

func (m *roomManager) GenCode() string {
	return fmt.Sprintf("%04d", len(m.rooms))
}
func (m *roomManager) AddConn(c *websocket.Conn, code string) error {
	room := m.get(code)
	if room == nil {
		return fmt.Errorf("room doesn't exist")
	}

	if err := room.Join(c); err != nil {
		return err
	}

	if room.IsFull() && !room.Started {
		room.Started = true
		go room.OnReady(room)
	}

	return nil
}
func (m *roomManager) RemoveConn(c *websocket.Conn, code string) error {
	room := m.get(code)
	if room == nil {
		return fmt.Errorf("room doesn't exist")
	}

	room.Leave(c)
	if room.IsEmpty() {
		m.RemoveRoom(code)
	}
	return nil
}

func (m *roomManager) CreateRoom(cfg string, onReady gameReadyCallabck) (code string) {
	var config gameConfig
	err := json.Unmarshal([]byte(cfg), &config)
	if err != nil {
		config = gameConfig{}
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	code = m.GenCode()
	m.rooms[code] = &gameRoom{
		Users: [2]user{
			{ch: make(chan *message)},
			{ch: make(chan *message)},
		},
		Config:  config,
		OnReady: onReady,
	}

	return
}
func (m *roomManager) RemoveRoom(code string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	close(m.rooms[code].Users[0].ch)
	close(m.rooms[code].Users[1].ch)
	delete(m.rooms, code)
}
