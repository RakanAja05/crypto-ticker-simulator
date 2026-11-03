package main

import (
	"log"

	"github.com/RakanAja05/crypto-ticker-simulator/gen/tickerpb"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		var tick tickerpb.Tick
		if err := proto.Unmarshal(msg, &tick); err != nil {
			log.Println("unmarshal:", err)
			continue
		}
		log.Printf("%s %.2f (%+.2f%%)", tick.Symbol, tick.Price, tick.ChangePct)
	}
}