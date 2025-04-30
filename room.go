package main

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

func newRoom(roomID string) *Room {
	return &Room{
		RoomID: roomID,
		Peers:  sync.Map{},
	}
}

func (room *Room) Join(myPeer *Peer) error {

	log.Printf(
		"%s join to room %s & peerID %s => \n", myPeer.Name, room.RoomID, myPeer.PeerID,
	)

	myPeerID := myPeer.PeerID

	payload, _ := json.Marshal(Payload{
		DestID: myPeerID,
		UserConnected: UserConnected{
			Name:    myPeer.Name,
			MyColor: myPeer.MyColor,
		},
	})

	return broadcast(room, myPeerID, Signal{
		Event:   JOIN_SIGNAL,
		Payload: string(payload),
	})
}

func (room *Room) Leave(myPeer *Peer) error {

	log.Printf(
		"%s leave to room %s => \n", myPeer.Name, room.RoomID,
	)

	payload, _ := json.Marshal(Payload{
		DestID: myPeer.PeerID,
		UserConnected: UserConnected{
			Name:    myPeer.Name,
			MyColor: myPeer.MyColor,
		},
	})

	return broadcast(room, myPeer.PeerID, Signal{
		Event:   LEAVE_SIGNAL,
		Payload: string(payload),
	})
}

func (room *Room) Close(myPeer *Peer) error {

	log.Printf(
		"%s close to room %s & peerID %s => \n", myPeer.Name, room.RoomID, myPeer.PeerID,
	)

	payload, _ := json.Marshal(Payload{
		DestID: myPeer.PeerID,
		UserConnected: UserConnected{
			Name:    myPeer.Name,
			MyColor: myPeer.MyColor,
		},
	})

	return broadcast(room, myPeer.PeerID, Signal{
		Event:   CLOSE_SIGNAL,
		Payload: string(payload),
	})
}

func broadcast(room *Room, myPeerID string, signal Signal) (err error) {
	room.Peers.Range(func(key, ref interface{}) bool {
		id := key.(string)

		if id != myPeerID && ref != nil {
			peer := ref.(*Peer)
			if err = peer.SafeWriteJSON(&signal); err != nil {
				log.Printf("failed to send signal. err: %v ", err)
				return false
			}
		}

		return true
	})

	return
}

func newPeer(
	peerID string, conn *websocket.Conn, name string, myColor string,
) *Peer {
	return &Peer{
		PeerID:  peerID,
		Conn:    conn,
		Name:    name,
		MyColor: myColor,
	}
}

func (p *Peer) SafeWriteJSON(v interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := p.Conn.WriteJSON(v); err != nil {
		return err
	}

	return nil
}
