package main

// WsServer : represent a chat server
// clients is a map for the registred client, register and unregister are chan to handle requests
type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

// newWsServer: create a new instance of WsServer
func newWsServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run : run our webserver and handle requests
func (server *WsServer) Run() {
	for {
		select {
		case client := <-server.register:
			server.registerClient(client)

		case client := <-server.unregister:
			server.unregisterClient(client)
		}
	}
}

// registerClient : register a client to the web server
func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, exists := server.clients[client]; exists {
		delete(server.clients, client)
	}
}
