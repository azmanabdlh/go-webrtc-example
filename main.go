package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	db *database
)

func init() {
	db = newDatabase()
}

func main() {

	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./public")))

	r.PathPrefix("/script.js").Handler(http.FileServer(http.Dir("./public")))

	r.HandleFunc("/ws/room/{room_id}", handle)

	log.Println("Server started at http://localhost:8000")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()
}

func handle(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	vars := mux.Vars(r)
	roomID := vars["room_id"]

	mySession := newSession(
		db,
		db.Find(roomID), // room object
	)

	// Handle WebSocket connection
	log.Println("Client connected:", conn.RemoteAddr())
	for {
		_, rawMsg, err := conn.ReadMessage()
		// send signal when client disconnected
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			mySession.Close()
			break
		}

		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		msg := Message{}
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			log.Print(err)
			log.Println("message: ", string(rawMsg))
			return
		}

		// log.Printf("Received message: %s\n", msg)
		log.Println("Received message:", msg.Event)
		err = mySession.Listen(conn, msg)
		if err != nil {
			log.Print("listen err => ", err)
			break
		}
	}
}
