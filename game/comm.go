package game

import (
	"bytes"
	"fmt"
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

func (c communicator) send(cmd string, args ...string) error {
	if c.conn == nil {
		return fmt.Errorf("connection is nil")
	}
	buf := bytes.Buffer{}

	buf.WriteString(cmd)
	for _, arg := range args {
		buf.WriteString("_" + arg)
	}

	return c.conn.WriteMessage(websocket.TextMessage, buf.Bytes())
}

func (c communicator) receive() (*message, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("connection is nil")
	}

	_, data, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return newMessage(data), nil
}
