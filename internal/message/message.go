package message

import (
	"encoding/json"
	"net"
)

type Delivery int

const (
	Identity Delivery = iota
	List
	Relay
)

type Message struct {
	Action    Delivery `json:"action"`
	Receivers []uint32 `json:"receivers,omitempty"`
	Body      string   `json:"body,omitempty"`
}

func (m *Message) DecodeIncomingMessage(conn net.Conn) error {

	decoder := json.NewDecoder(conn)

	err := decoder.Decode(m)

	if err != nil {
		return err
	}

	return nil

}
