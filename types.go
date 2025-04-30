package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

const (
	JOIN_SIGNAL          string = "JOIN_SIGNAL"
	ANSWER_SIGNAL               = "ANSWER_SIGNAL"
	OFFER_SIGNAL                = "OFFER_SIGNAL"
	LEAVE_SIGNAL                = "LEAVE_SIGNAL"
	CLOSE_SIGNAL                = "CLOSE_SIGNAL"
	ICE_CANDIDATE_SIGNAL        = "ICE_CANDIDATE_SIGNAL"

	JOINED_SIGNAL = "JOINED_SIGNAL"

	SEND_ANSWER_SIGNAL        = "SEND_ANSWER_SIGNAL"
	SEND_OFFER_SIGNAL         = "SEND_OFFER_SIGNAL"
	NEW_PARTICIPANT_SIGNAL    = "NEW_PARTICIPANT_SIGNAL"
	SEND_ICE_CANDIDATE_SIGNAL = "SEND_ICE_CANDIDATE_SIGNAL"
)

type Room struct {
	RoomID string
	// Peers  map[string]*Peer

	// for safe using lock
	Peers sync.Map
}

type SessionManager struct {
	room *Room
	db   *database

	myPeerID string
}

type Peer struct {
	// unique identifier for the peer
	PeerID  string
	Conn    *websocket.Conn
	Name    string
	MyColor string // my video background color

	mu sync.Mutex
}

type UserConnected struct {
	Name    string `json:"name"`
	MyColor string `json:"myColor"`
}

type Payload struct {
	Type          string        `json:"type"`
	Sdp           string        `json:"sdp"`
	OriginID      string        `json:"originID"`
	DestID        string        `json:"destID"`
	UserConnected UserConnected `json:"userConnected"`
}

type Candidate struct {
	Candidate     string `json:"candidate"`
	SdpMLineIndex int    `json:"sdpMLineIndex"`
	SdpMid        string `json:"sdpMid"`
}

type CandidatePlain struct {
	Candidate string `json:"candidate"`
	OriginID  string `json:"originID"`
	DestID    string `json:"destID"`
}

type Message struct {
	Event   string `json:"event"`
	Payload string `json:"payload"`
}

type database struct {
	memory sync.Map
}

type Signal struct {
	Event   string `json:"event"`
	Payload string `json:"payload"`
}
