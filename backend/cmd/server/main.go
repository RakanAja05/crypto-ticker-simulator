package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/RakanAja05/crypto-ticker-simulator/gen/tickerpb"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // dev only
}

func handleTicker(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer conn.Close()

	prices := map[string]float64{"BTC": 67000, "ETH": 3500, "SOL": 150}
	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

	for range t.C {
		for sym, p := range prices {
			change := (rand.Float64() - 0.5) * 0.01 // random walk ±0.5%
			newP := p * (1 + change)
			prices[sym] = newP

			data, err := proto.Marshal(&tickerpb.Tick{
				Symbol:    sym,
				Price:     newP,
				ChangePct: change * 100,
				Timestamp: time.Now().UnixMilli(),
			})
			if err != nil {
				log.Println("marshal:", err)
				return
			}
			if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
				return // client disconnect
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleTicker)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}