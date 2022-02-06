package game

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kananb/chess"
)

type player struct {
	Conn *websocket.Conn
	ch   chan *message
}

func (p player) IsConnected() bool {
	return p.Conn != nil
}

type gameReadyCallabck func(*gameRoom)

type gameRoom struct {
	players [2]*player
	Board   *chess.Board
	OnReady gameReadyCallabck
	Started bool
}

func (r *gameRoom) Join(p *player) error {
	var i int
	if r.players[0] == nil {
		i = 0
	} else if r.players[1] == nil {
		i = 1
	} else {
		return fmt.Errorf("room is full")
	}

	p.ch = make(chan *message)
	r.players[i] = p

	return nil
}
func (r *gameRoom) Leave(p *player) {
	var i int
	if r.players[0] == p {
		i = 0
	} else if r.players[1] == p {
		i = 1
	} else {
		return
	}

	close(r.players[i].ch)
	r.players[i] = nil
}

func (r *gameRoom) Broadcast(cmd string, args ...string) {
	comms := []communicator{
		{r.players[0].Conn},
		{r.players[1].Conn},
	}

	for _, comm := range comms {
		comm.send(cmd, args...)
	}
}

func (r *gameRoom) IsEmpty() bool {
	return r.players[0] == nil && r.players[1] == nil
}
func (r *gameRoom) IsFull() bool {
	return r.players[0] != nil && r.players[1] != nil
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
func (m *roomManager) AddPlayer(p *player, code string) error {
	room := m.get(code)
	if room == nil {
		return fmt.Errorf("room doesn't exist")
	}

	if err := room.Join(p); err != nil {
		return err
	}

	if room.IsFull() && !room.Started {
		room.Started = true
		go room.OnReady(room)
	}

	return nil
}
func (m *roomManager) RemovePlayer(p *player, code string) error {
	if p == nil {
		return fmt.Errorf("player is nil")
	}

	room := m.get(code)
	if room == nil {
		return fmt.Errorf("room doesn't exist")
	}

	room.Leave(p)
	if room.IsEmpty() {
		m.RemoveRoom(code)
	}
	return nil
}

func (m *roomManager) CreateRoom(onReady gameReadyCallabck) (code string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	code = m.GenCode()
	m.rooms[code] = &gameRoom{
		OnReady: onReady,
	}

	return
}
func (m *roomManager) RemoveRoom(code string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.rooms, code)
}
