package game

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/kananb/chess"
)

type player struct {
	Conn *websocket.Conn
	Side chess.SideColor
}

func (p player) IsConnected() bool {
	return p.Conn != nil
}

type gameRoom struct {
	P1, P2 player
	Board  *chess.Board
}

func (r *gameRoom) Join(conn *websocket.Conn) error {
	if !r.P1.IsConnected() {
		r.P1.Conn = conn
	} else if !r.P2.IsConnected() {
		r.P2.Conn = conn
	} else {
		return fmt.Errorf("room is full")
	}

	return nil
}
func (r *gameRoom) Leave(conn *websocket.Conn) {
	if r.P1.Conn == conn {
		r.P1.Conn = nil
	} else if r.P2.Conn == conn {
		r.P2.Conn = nil
	}
}

type roomManager struct {
	rooms map[string]*gameRoom
	lock  sync.RWMutex
}

func (m *roomManager) get(code string) *gameRoom {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.rooms[code]
}
func (m *roomManager) set(code string, room *gameRoom) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.rooms[code] = room
}

func (m *roomManager) GenCode() string {
	return fmt.Sprintf("%04d", len(m.rooms))
}
func (m *roomManager) AddConn(conn *websocket.Conn, code string) error {
	room := m.get(code)
	if room == nil {
		return fmt.Errorf("room doesn't exist")
	}

	if err := room.Join(conn); err != nil {
		return err
	}

	return nil
}
func (m *roomManager) RemoveConn(conn *websocket.Conn, code string) error {
	room := m.get(code)
	if room == nil {
		return fmt.Errorf("room doesn't exist")
	}

	room.Leave(conn)
	return nil
}

func (m *roomManager) CreateRoom() (code string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	code = m.GenCode()
	m.rooms[code] = &gameRoom{Board: chess.StartingPosition()}

	return
}
func (m *roomManager) RemoveRoom(code string) {
	m.set(code, nil)
}
