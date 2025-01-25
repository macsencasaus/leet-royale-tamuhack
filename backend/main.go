package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {

	mux := http.NewServeMux()

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			panic("Error creating websocket")
		}

		conn.WriteMessage(websocket.TextMessage, []byte("Poop fart"))
	})

	fmt.Println("Serving on 10.246.176.24:6969")
	http.ListenAndServe("10.246.176.24:6969", mux)
}
