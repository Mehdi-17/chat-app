package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// writeWait : when the server send a message to a client, it expects the write operation to complete within writeWait limit
	writeWait = 10 * time.Second

	// pongWait : max time that the server waits for a pong response. Server waits for a pong to ensure the connection is alive
	pongWait = 60 * time.Second

	// pingPeriod : ping interval to ensure that the connection is active, less than pongWait to allow time for the client to respond
	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 10000
)

// upgrader : Define a variable to set the read and write buffer size of my chat app
// A large buffer help to handle larger message and avoid buffer resizing but take more memory
var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

// Client : Defines the Client type to handle the websocket connection for a client and an instance of the web server
type Client struct {
	conn     *websocket.Conn
	WsServer *WsServer
}

func newClient(conn *websocket.Conn, server *WsServer) *Client {
	return &Client{
		conn:     conn,
		WsServer: server,
	}
}

// ServeWs : upgrade a http request to a Websocket connection and create a client for the ws conn
func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {
	// security risk but, accept websocket connection from any origins
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	//upgrade the http connection to a Websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	//create a client for the new websocket connection
	client := newClient(conn, wsServer)

	//TODO : to implement writePump
	//go client.writePump()
	go client.readPump()

	wsServer.register <- client

	fmt.Printf("A new client has joined us => %v\n", client)
}

// readPump : handler to read message from the websocket connection
func (client *Client) readPump() {
	defer func() {
		//todo : call disconnect func
		//client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))

	//when receive a pong, updates the read deadline to extend the allowed time for the next Pong message
	client.conn.SetPongHandler(func(appData string) error {
		err := client.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return nil
	})

	for {
		_, jsonMessage, err := client.conn.ReadMessage()

		if err != nil {
			log.Printf("Unexpected error during readMessage : %v", err)
			break
		}

		client.WsServer.broadcast <- jsonMessage
	}
}
