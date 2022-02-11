package data

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/kananb/chess"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type Player struct {
	ID    string
	Clock int
}

// Stores information about a chess game.
// Players[0] corresponds to the white side and Players[1] corresponds to black
type GameState struct {
	Players        [2]Player
	FEN            string
	SideToMove     chess.SideColor
	History        []string
	TimeOfLastMove int64
	Increment      int
	Result         string
}
type userGameData struct {
	Color      chess.SideColor
	WhiteClock int
	BlackClock int
	FEN        string
	SideToMove chess.SideColor
	History    []string
	Result     string
}

func NewGameState() *GameState {
	return &GameState{
		SideToMove: chess.White,
		History:    []string{},
	}
}

func (s *GameState) Marshal(playerID string) string {
	color := chess.White
	if s.Players[1].ID == playerID {
		color = chess.Black
	}
	u := &userGameData{
		Color:      color,
		WhiteClock: s.Players[0].Clock,
		BlackClock: s.Players[1].Clock,
		FEN:        s.FEN,
		SideToMove: s.SideToMove,
		History:    s.History,
		Result:     s.Result,
	}

	val, err := json.Marshal(u)
	if err != nil {
		panic("failed to marshall userGameState")
	}

	return string(val)
}
func (s *GameState) String() string {
	val, err := json.Marshal(s)
	if err != nil {
		panic("failed to marshall GameState")
	}

	return string(val)
}

type StateManager interface {
	genKey() string

	Get(key string) *GameState
	Set(key string, state *GameState) bool
	Del(key string) bool

	Sub(channel string) chan string
	Unsub(channel string, in chan string) bool
	Pub(msg, channel string, in chan string) bool
}

type GameConfig struct {
	Duration  int
	Increment int
	PlayAs    chess.SideColor
}

func NewGameConfig(cfg string) *GameConfig {
	config := new(GameConfig)
	err := json.Unmarshal([]byte(cfg), config)
	if err != nil {
		return nil
	}

	return config
}

type ChessManager struct {
	manager StateManager
}

func NewChessManager(manager StateManager) ChessManager {
	return ChessManager{manager}
}

func (m ChessManager) Join(code, playerID string) (id string, in chan string, err error) {
	state := m.manager.Get(code)
	i := -1
	if state == nil {
		return "", nil, fmt.Errorf("room doesn't exist")
	} else if state.Players[0].ID == playerID {
		i = 0
	} else if state.Players[1].ID == playerID {
		i = 1
	}

	if i == -1 {
		return "", nil, fmt.Errorf("room is full")
	} else if state.Players[i].ID == "" {
		state.Players[i].ID = fmt.Sprintf("%08X", rand.Int31())
		if state.FEN == "" {
			state.FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
			state.TimeOfLastMove = time.Now().UnixMilli()
		}
		m.manager.Set(code, state)
		m.manager.Pub("NEW", code, nil)
	}

	return state.Players[i].ID, m.manager.Sub(code), nil
}
func (m ChessManager) Create(config GameConfig) (code, id string, in chan string) {
	code = m.manager.genKey()
	in = m.manager.Sub(code)

	state := m.manager.Get(code)
	if state == nil {
		m.manager.Unsub(code, in)
		return "", "", nil
	}

	if config.Duration <= 0 {
		config.Duration = 59990
	}
	state.Players[0].Clock = config.Duration + 9
	state.Players[1].Clock = config.Duration + 9
	if config.Increment < 0 {
		config.Increment = 0
	}
	state.Increment = config.Increment

	i := 0
	if config.PlayAs == chess.White {
		state.Players[0].ID = fmt.Sprintf("%08X", rand.Int31())
	} else if config.PlayAs == chess.Black {
		state.Players[1].ID = fmt.Sprintf("%08X", rand.Int31())
		i = 1
	} else {
		i = rand.Intn(2)
		state.Players[i].ID = fmt.Sprintf("%08X", rand.Int31())
	}
	m.manager.Set(code, state)

	return code, state.Players[i].ID, in
}
func (m ChessManager) Leave(code string, in chan string) {
	m.manager.Unsub(code, in)
}

func (m ChessManager) Get(code string) *GameState {
	return m.manager.Get(code)
}
func (m ChessManager) Set(code string, state *GameState) bool {
	return m.manager.Set(code, state)
}
func (m ChessManager) Notify(state *GameState, code string, in chan string) bool {
	if ok := m.manager.Set(code, state); !ok {
		return false
	}
	return m.manager.Pub("NEW", code, in)
}
