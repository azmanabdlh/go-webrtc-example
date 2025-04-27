package main

import "sync"

func newDatabase() *database {
	return &database{
		memory: sync.Map{},
	}
}

func (db *database) Find(roomID string) *Room {
	if ref, ok := db.memory.Load(roomID); ok {
		room := ref.(*Room)
		return room
	}

	// create room
	room := newRoom(roomID)
	db.memory.Store(roomID, room)
	return room
}
