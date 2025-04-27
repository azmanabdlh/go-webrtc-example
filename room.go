package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func newRoom(roomID string) *Room {
	return &Room{
		RoomID: roomID,
		Peers:  make(map[string]*Peer),
	}
}

func (room *Room) Join(myPeer *Peer) {

	log.Printf(
		"%s join to room %s & peerID %s => \n", myPeer.Name, room.RoomID, myPeer.PeerID,
	)

	myPeerID := myPeer.PeerID
	room.Peers[myPeerID] = myPeer

	log.Println(
		"rooms peers => ", len(room.Peers),
	)

	payload, _ := json.Marshal(Payload{
		DestID: myPeerID,
		UserConnected: UserConnected{
			Name:    myPeer.Name,
			MyColor: myPeer.MyColor,
		},
	})

	broadcast(room, myPeerID, Signal{
		Event:   JOIN_SIGNAL,
		Payload: string(payload),
	})
}

func (room *Room) Leave(myPeer *Peer) {

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

	broadcast(room, myPeer.PeerID, Signal{
		Event:   LEAVE_SIGNAL,
		Payload: string(payload),
	})
}

func (room *Room) Close(myPeer *Peer) {

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

	broadcast(room, myPeer.PeerID, Signal{
		Event:   CLOSE_SIGNAL,
		Payload: string(payload),
	})
}

func broadcast(room *Room, myPeerID string, signal Signal) {
	for id, peer := range room.Peers {
		if id == myPeerID {
			continue
		}

		if err := peer.Conn.WriteJSON(&signal); err != nil {
			log.Printf("failed to send signal. err: %v ", err)
			break
		}
	}
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
