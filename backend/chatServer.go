package main

// WsServer : represent a chat server
// clients is a map for the registred client, register and unregister are chan to handle requests
type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

//TODO: implement features
