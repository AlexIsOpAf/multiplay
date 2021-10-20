package hub

import (
	"log"
	"net"

	"github.com/cornelk/hashmap"
	uuid "github.com/nu7hatch/gouuid"
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
------------------------------------------------------------------------------
*/

// Hub allows for no hard limit to accepted client, implementation should allow for threadsafe assignments
type Hub struct {

	// Created Server
	ls net.Listener

	connectedClients *hashmap.HashMap

	// Hub needs to have tunneled access from clients - sockets?
}

func ServerConn(port, conn_type string) error {

	lis, err := net.Listen(conn_type, port)

	if err != nil {
		return err
	}

	newHub := &Hub{
		ls:               lis,
		connectedClients: &hashmap.HashMap{},
	}

	// No auth needed: For every connection, accept incoming
	for {
		conn, err := lis.Accept()

		if err != nil {
			return err
		}

		// For every accepted incoming, handle in new goroutine
		go newHub.handleConnection(conn)
	}
}

// For every accepted incoming, allocate
func (hub *Hub) handleConnection(conn net.Conn) {

	// Assign generated unique ID from external module to conn - don't need to manage my own
	// generates 16 bytes rather than 8 from uint64
	u, err := uuid.NewV4()

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	hub.connectedClients.Set(u.String(), conn)

	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err = conn.Read(buf)

	if err != nil {
		log.Fatal("Error reading:", err.Error())
		return
	}

	conn.Write([]byte("Message received."))
}

func (hub *Hub) Close() error {
	err := hub.ls.Close()

	if err != nil {
		return err
	}

	return nil
}
