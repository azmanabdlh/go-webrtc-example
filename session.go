package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

func generateRandHash(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, length)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func newSession(db *database, room *Room) *SessionManager {
	return &SessionManager{room, db}
}

func (sm *SessionManager) Listen(conn *websocket.Conn, msg Message) error {

	switch msg.Event {
	case NEW_PARTICIPANT_SIGNAL:
		myPeerID := generateRandHash(8) // my uniqueID
		payload := struct {
			Name    string `json:"name"`
			RoomID  string `json:"roomID"`
			MyColor string `json:"myColor"`
		}{}

		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			return err
		}

		myPeer := newPeer(myPeerID, conn, payload.Name, payload.MyColor)

		payloadByte, _ := json.Marshal(Payload{
			OriginID: myPeerID,
		})

		if err := conn.WriteJSON(&Signal{
			Event:   JOINED_SIGNAL,
			Payload: string(payloadByte),
		}); err != nil {
			return err
		}

		sm.room.Join(myPeer)

	case SEND_ICE_CANDIDATE_SIGNAL:
		payload := CandidatePlain{}
		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			return err
		}

		room := db.Find(
			sm.room.RoomID,
		)

		if myPeer, ok := room.Peers[payload.DestID]; ok {
			err := myPeer.Conn.WriteJSON(&Message{
				Event:   ICE_CANDIDATE_SIGNAL,
				Payload: msg.Payload,
			})
			if err != nil {
				return err
			}
		}

		return fmt.Errorf("peerID %s not found for roomID %s", payload.DestID, room.RoomID)

	case SEND_ANSWER_SIGNAL:
		fallthrough
	case SEND_OFFER_SIGNAL:
		payload := Payload{}
		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			return err
		}

		destID := payload.DestID

		room := db.Find(
			sm.room.RoomID,
		)

		signal := OFFER_SIGNAL
		if payload.Type == "answer" {
			signal = ANSWER_SIGNAL
		}

		log.Println("payload.UserConnected => ", payload.UserConnected)

		payloadStr, _ := json.Marshal(payload)

		if myPeer, ok := room.Peers[destID]; ok {
			if err := myPeer.Conn.WriteJSON(&Message{
				Event:   signal,
				Payload: string(payloadStr),
			}); err != nil {
				return err
			}

			return nil
		}

		return fmt.Errorf("peerID %s not found for roomID %s", destID, room.RoomID)

	}
	return nil
}
