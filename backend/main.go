package main

import (
	"flag"
	"log"
	"net/http"
)

// addr: configuration options used to specify to the http server where the Ws server will listen connections
// Can be overridden by a command-line args, e.g :./main.go -addr :9090
var addr = flag.String("addr", ":8080", "http server address")

func main() {
	//Parse : reads the command-line args and update the config of the flag. If we doesn't provide args, it will take our default value ":8080"
	flag.Parse()

	wsServer := newWsServer()
	go wsServer.Run()

	//HandleFunc: handler for http request with the path '/ws'. The http request will be upgrade to a websocket conn and we will create a client
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(wsServer, w, r)
	})

	//ListenAndServe : is used to listen to the port we specify on addr
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
