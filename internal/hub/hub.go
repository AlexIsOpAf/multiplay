package hub

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/cornelk/hashmap"
	"github.com/google/uuid"
	"github.com/multiplay/internal/message"
)

/*
------------------------------------------------------------------------------
Information/Assumptions surrounding HUB:

1. Communication between nodes in the system must occur at the network layer - TCP server/80
2. Clients route through hub to deliver messages to other connected clients
3. No auth for connected clients
4. No indication there will/should be multiple hubs present
5. Connected clients are given unique ID (uint64)
6. No hard limit to clients that can be connected to 1 hub
7. Assume you can connect external clients to hub i.e postman for test cases
8. userID => uint64?
9. No indication of scale hub has to adhere to
------------------------------------------------------------------------------
*/

// Hub allows for no hard limit to accepted client, implementation should allow for threadsafe assignments
type Hub struct {

	// Created Server
	ls net.Listener

	// thread safe hashmap
	connectedClients *hashmap.HashMap
}

func NewHub(port, conn_type string) (*Hub, error) {
	lis, err := net.Listen(conn_type, port)

	if err != nil {
		return nil, err
	}

	newHub := &Hub{
		ls:               lis,
		connectedClients: &hashmap.HashMap{},
	}

	return newHub, nil
}

func (h *Hub) ServerConn() error {

	// No auth needed: For every connection, accept incoming
	for {
		conn, err := h.ls.Accept()

		if err != nil {
			return err
		}

		// For every accepted incoming, handle in new goroutine
		go h.HandleConnection(conn)
	}
}

// For every accepted incoming
func (hub *Hub) HandleConnection(conn net.Conn) {

	// Assign generated unique ID from external module to conn - don't need to manage my own
	// Can generate UUID and use it's ID to get uint32
	u, err := uuid.NewUUID()

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// set uuid's ID to the current client connection being handled
	hub.connectedClients.Set(u.ID(), conn)

	// Main conn process loop

	for {

		var messageContainer message.Message

		err := messageContainer.DecodeIncomingMessage(conn)

		if err != nil {
			log.Fatal("Error: ", err.Error())
			return
		}

		switch messageContainer.Action {

		case message.Identity:
			conn.Write([]byte(ConvertUintToString(u.ID())))

		case message.List:
			keyValue := hub.connectedClients.Iter()
			var clientList []uint32

			for k := range keyValue {
				id := k.Key.(uint32)
				if id != u.ID() {
					clientList = append(clientList, id)
				}
			}
			returnMessage := FormatListMessage(clientList)
			conn.Write([]byte(returnMessage))

		case message.Relay:
			// Omit error handling as we are only told to deliver the message

			if len(messageContainer.Body) > 1048576 {
				continue
			}

			if len(messageContainer.Receivers) < 256 {
				for _, v := range messageContainer.Receivers {
					if val, ok := hub.connectedClients.Get(v); ok {
						val.(net.Conn).Write([]byte(messageContainer.Body))
					}
				}
			}
		}
	}
}

func FormatListMessage(listOfClients []uint32) string {

	var builder strings.Builder

	// Title
	builder.WriteString("Connected Clients:\n")

	for _, v := range listOfClients {
		fmt.Fprintf(&builder, ConvertUintToString(v)+"\n")
	}

	return builder.String()
}

func ConvertUintToString(val uint32) string {
	return fmt.Sprint(val)
}
