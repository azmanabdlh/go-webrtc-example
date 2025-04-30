# Journey of Learning WebRTC with Go & JavaScript

This repository documents my **journey of learning WebRTC** using a Go backend (Pion WebRTC, Gorilla WebSocket, mux) and HTML/JavaScript frontend.  
The main goal is to understand signaling concepts, relay, and implement a simple video conference from scratch.

---

## Features

- Signaling server with Go (WebSocket)
- Simple room & peer management
- Video relay between browsers (SFU-lite)
- Easy-to-understand HTML/JS client code

---
## How to Flow Run

1. **Run the Go server:**
   ```sh
   make run
   ```
   The server will run at [http://localhost:8000](http://localhost:8000)

2. **Open your browser** at `http://localhost:8000/`
3. **Join a room** with your name & room ID (you can share the room ID with friends)
4. **Try opening from another tab/browser** to simulate multi-user

### Example:
  1. **User A join the room**
  - Enters name & room ID.
  - A send signal `NEW_PARTICIPANT_SIGNAL` to the websocket server
  - websocket server no relay yet (no other peers).
  
  2. **User B joins the same room**
  - Enters name & room ID.
  - B send signal `NEW_PARTICIPANT_SIGNAL` to the websocket server

  3. **Server detect User A is present and send `JOIN_SIGNAL` to User A with User B data:**
  ```json
  {
    "name": "B-name",
    "myColor": "B-myColor",
    "destID": "A-ID"
  }
  ```
  
  4. **User A got responder from server relay `JOIN_SIGNAL`**
  - User A create a WebRTC connection for User B using User B data
  - User A add remote video User B to DOM
  - User A send offer signal `SEND_OFFER_SIGNAL` to websocket server with User A data
   ```json
  {
    "name": "A-name",
    "myColor": "A-myColor",
    "originID": "A-ID",
    "destID": "B-ID"
  }
  ```

  5. **websocket server relay the offer signal `OFFER_SIGNAL` to User B**
  - User B create a WebRTC connection for User A using User A data
  - User B add remote video User A to DOM
  - User B send answer signal `SEND_ANSWER_SIGNAL` to User A

  6. **websocket server relay the answer signal `ANSWER_SIGNAL` to User A**
  - User A got answer and exchange ICE candidates connection

  7. done

---

## Dependencies

- [Pion WebRTC](https://github.com/pion/webrtc)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)
- [Gorilla Mux](https://github.com/gorilla/mux)

---

## License

MIT

---

> _This repository was created as personal documentation and practice.  
> Hopefully, it will be useful for friends who also want to learn WebRTC from scratch!_