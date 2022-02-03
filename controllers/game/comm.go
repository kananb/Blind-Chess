package game

import (
	"bytes"
	"strings"

	"github.com/gorilla/websocket"
)

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
	} else {
		msg.Args = []string{""}
	}

	return msg
}

type communicator struct {
	conn *websocket.Conn
}

func (c communicator) send(code string, args ...string) {
	buf := bytes.Buffer{}

	buf.WriteString(code)
	for _, arg := range args {
		buf.WriteString("_" + arg)
	}

	c.conn.WriteMessage(websocket.TextMessage, buf.Bytes())
}

func (c communicator) receive() (*message, error) {
	_, data, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return newMessage(data), nil
}
