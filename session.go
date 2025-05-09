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
	return &SessionManager{room, db, ""}
}

func (sm *SessionManager) Listen(conn *websocket.Conn, msg Message) error {

	switch msg.Event {
	case NEW_PARTICIPANT_SIGNAL:
		myPeerID := generateRandHash(8) // my uniqueID
		sm.myPeerID = myPeerID
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

		sm.room.Peers.Store(myPeerID, myPeer)

		sm.room.Join(myPeer)

	case SEND_ICE_CANDIDATE_SIGNAL:
		payload := CandidatePlain{}
		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			return err
		}

		room := db.Find(
			sm.room.RoomID,
		)

		if ref, ok := room.Peers.Load(payload.DestID); ok && ref != nil {
			myPeer := ref.(*Peer)
			err := myPeer.Conn.WriteJSON(&Message{
				Event:   ICE_CANDIDATE_SIGNAL,
				Payload: msg.Payload,
			})
			if err != nil {
				return err
			}

			return nil
		}

		return fmt.Errorf("peerID %s not SEND_ICE_CANDIDATE_SIGNAL found for roomID %s", payload.DestID, room.RoomID)

	case SEND_ANSWER_SIGNAL:
		fallthrough
	case SEND_OFFER_SIGNAL:
		payload := Payload{}
		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			return err
		}

		peerID := payload.DestID

		room := db.Find(
			sm.room.RoomID,
		)

		signal := OFFER_SIGNAL
		if payload.Type == "answer" {
			signal = ANSWER_SIGNAL
		}

		payloadStr, _ := json.Marshal(payload)

		if ref, ok := room.Peers.Load(peerID); ok && ref != nil {
			myPeer := ref.(*Peer)
			if err := myPeer.Conn.WriteJSON(&Message{
				Event:   signal,
				Payload: string(payloadStr),
			}); err != nil {
				return err
			}

			return nil
		}

		return fmt.Errorf("peerID %s not found for roomID %s", peerID, room.RoomID)

	}
	return nil
}

func (sm *SessionManager) Close() {
	key := sm.myPeerID

	if ref, ok := sm.room.Peers.Load(key); ok && ref != nil {
		myPeer := ref.(*Peer)
		err := sm.room.Leave(myPeer)
		if err != nil {
			log.Printf("Error leaving room %s, peer %s:\n", err, myPeer.PeerID)
			return
		}

		sm.room.Peers.Delete(key)
	}
}
