package main

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"net/http"
)

func main() {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(conn socketio.Conn) error {
		conn.SetContext("")
		fmt.Println("Connected :", conn.ID())
		return nil
	})

	server.OnEvent("/chat", "msg", func(conn socketio.Conn, msg string) string {
		conn.SetContext(msg)
		fmt.Println("Message receive : ", conn.ID(), msg)
		//TODO : implement handler
		return "recv " + msg
	})

	server.OnDisconnect("/", func(conn socketio.Conn, reason string) {
		fmt.Println("Diconnected : ", conn.ID(), reason)
	})

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal("error during serve ", err)
		}
	}()

	defer func(server *socketio.Server) {
		err := server.Close()
		if err != nil {
			log.Fatal("Disconnection error ", err)
		}
	}(server)

	http.Handle("/socket.io/", server)

	err := http.ListenAndServe(":8080", nil)
	fmt.Println("Serving on localhost:8080")

	if err != nil {
		log.Fatal("Server error:", err)
	}

}
