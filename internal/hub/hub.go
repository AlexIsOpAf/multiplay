package hub

import (
	"fmt"
	"log"
	"net"

	"github.com/cornelk/hashmap"
	"github.com/google/uuid"
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

	// list of clients
	listOfClients []net.Conn
}

func ServerConn(port, conn_type string) error {

	lis, err := net.Listen(conn_type, port)

	if err != nil {
		return err
	}

	var base []net.Conn

	newHub := &Hub{
		ls:               lis,
		connectedClients: &hashmap.HashMap{},
		listOfClients:    base,
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
	// Can generate UUID and use it's ID to get uint32
	u, err := uuid.NewUUID()

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	hub.connectedClients.Set(conn.RemoteAddr().String(), u.ID())
	hub.listOfClients = append(hub.listOfClients, conn)

	// for {
	// 	data, err := bufio.NewReader(conn).ReadString('\n')

	// 	if err != nil {
	// 		log.Fatal("Error: ", err.Error())
	// 		return
	// 	}

	// 	// "Endpoints"
	// 	if data == "Who Am I?\n" {
	// 		conn.Write([]byte(hub.ConvertUintToString(u.ID())))
	// 	}

	// 	if data == "Who is here?\n" {
	// 		// Iterate through entire list of clients connected
	// 		for i := 0; i < len(hub.listOfClients); i++ {
	// 			//
	// 			if val, ok := hub.connectedClients.GetStringKey(hub.listOfClients[i].RemoteAddr().String()); ok {
	// 				if val.(uint32) != u.ID() {
	// 					conn.Write([]byte(hub.ConvertUintToString(val.(uint32))))
	// 				}
	// 			}
	// 		}
	// 	}

	// }
}

func (hub *Hub) ConvertUintToString(val uint32) string {
	return fmt.Sprint(val)
}
