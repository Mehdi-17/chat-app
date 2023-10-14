package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// upgrader : Define a variable to set the read and write buffer size of my chat app
// A large buffer help to handle larger message and avoid buffer resizing but take more memory
var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

// Client : Defines the Client type to handle the websocket connection for a client
type Client struct {
	conn *websocket.Conn
}

func newClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

// ServeWs : upgrade a http request to a Websocket connection and create a client for the ws conn
func ServeWs(w http.ResponseWriter, r *http.Request) {
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
	client := newClient(conn)

	fmt.Printf("A new client has joined us => %v\n", client)
}
