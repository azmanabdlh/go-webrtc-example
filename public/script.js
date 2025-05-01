
const ANSWER_SIGNAL = "ANSWER_SIGNAL";
const OFFER_SIGNAL = "OFFER_SIGNAL";    
const JOIN_SIGNAL = "JOIN_SIGNAL";
const JOINED_SIGNAL = "JOINED_SIGNAL";
const ICE_CANDIDATE_SIGNAL = "ICE_CANDIDATE_SIGNAL";
const LEAVE_SIGNAL = "LEAVE_SIGNAL";

const NEW_PARTICIPANT_SIGNAL = "NEW_PARTICIPANT_SIGNAL";
const SEND_OFFER_SIGNAL = "SEND_OFFER_SIGNAL";
const SEND_ANSWER_SIGNAL = "SEND_ANSWER_SIGNAL";
const SEND_ICE_CANDIDATE_SIGNAL = "SEND_ICE_CANDIDATE_SIGNAL";    


let myStreamCam = new MediaStream();

const containerEl = document.querySelector('.container .main');
const form = document.querySelector('.chat_window__writer form');

const chatWindowBtn = document.querySelector('.chat_window__btn button');
const chatWindow = document.querySelector('.chat_window');
const chatBox = document.querySelector('.chat_window__chatbox');


let ws;
let streamChannel;
let yourname = "";
let roomID = "";
let myLocalPeerID = "";    
const myConns = new Object();
const myLocalColor = generateRandomColor();



chatWindowBtn.addEventListener('click', function(e) {
  let btnText = "Close chat";

  e.target.parentElement.classList.toggle("active");  
  chatWindow.classList.toggle("active");

  if (chatWindow.classList.contains("active")) {
  	btnText = "Close chat";
  } else {
  	btnText = "Open chat";
  }

  e.target.innerText = btnText;
});

form.addEventListener('submit', function(e) {
  e.preventDefault();

  const message = e.target.elements[0].value;
  
  let data = {
    name: yourname,
    message: message,
    color: myLocalColor,
  }

  form.reset();
  streamChannel.send(JSON.stringify(data));
  addChatMessage(chatBox, data);
});

window.onload = function() {

  const host = window.location.host;
  const protocol = window.location.protocol;
  let wsProtocol = "ws";

  if (protocol == "https:") {
    wsProtocol = "wss";
  }

  ws = new WebSocket(wsProtocol + '://'+ host +'/ws/room/' + roomID);

  ws.onopen = main;
  ws.onmessage = handleOnMessage;
  ws.onclose = function() {
    console.log("WebSocket connection closed");
  };
  ws.onerror = function() {
    console.log("WebSocket connection error");
  };
};

window.onclose = function () {
  if (ws != null) {
    ws.close();
  }
}


function generateRandomColor() {
  const letters = '0123456789ABCDEF';
  let color = '#';
  for (let i = 0; i < 6; i++) {
    color += letters[Math.floor(Math.random() * 16)];
  }
  return color;
}


function addChatMessage(container, object) {
  const chatEl1 = document.createElement('ul');
  const chatEl2 = document.createElement('li');

  const chatNameEl = document.createElement('p');
  const chatMessageEl = document.createElement('p');

  const date = new Date();
  
  // human date only show time
  const humanDate = date.toLocaleTimeString('id-ID', {
    hour: '2-digit',
    minute: '2-digit',    
  });

  chatNameEl.classList.add('chatbox__user_name');
  chatMessageEl.classList.add('chatbox__user_message');

  chatNameEl.style.color = object.color;
  chatNameEl.innerText = object.name + " - Pukul " + humanDate;
  chatMessageEl.innerText = object.message;


  chatEl2.appendChild(chatNameEl);
  chatEl2.appendChild(chatMessageEl);

  chatEl1.appendChild(chatEl2)
  container.appendChild(
   chatEl1
  );
}

function addCam(myID, name, myColor, stream, container) {
  const videoCam = document.createElement('video');
  videoCam.autoplay = true;
  videoCam.muted = true;
  videoCam.srcObject = stream;

  const label = document.createElement('div');
  label.className = 'name';
  label.textContent = name;

  const box = document.createElement('div');
  box.classList.add('video_cam');  
  box.setAttribute('data-id', myID);
  box.style.border = `3px solid ${myColor}`;

  box.appendChild(videoCam);
  box.appendChild(label);
  container.appendChild(box);

  syncVideoLayout(container)
}

function syncVideoLayout(grid) {  
  const gridLen = grid.children.length;
  let cols = 1, rows = 1;

  if (gridLen === 1) {
    cols = 1; rows = 1;
  } else if (gridLen === 2) {
    cols = 2; rows = 1;
  } else if (gridLen <= 4) {
    cols = 2; rows = 2;
  } else if (gridLen <= 6) {
    cols = 3; rows = 2;
  } else if (gridLen <= 9) {
    cols = 3; rows = 3;
  } else if (gridLen <= 12) {
    cols = 4; rows = 3;
  } else {
    cols = Math.ceil(Math.sqrt(count));
    rows = Math.ceil(count / cols);
  }
  grid.style.gridTemplateColumns = `repeat(${cols}, 1fr)`;
  grid.style.gridTemplateRows = `repeat(${rows}, 1fr)`;
}


function removeCam(myID, container) {
  const videoCam = container.querySelector(`.video_cam[data-id="${myID}"]`);
  console.log("remove cam => ", videoCam);
  console.log("remove id => ", myID);
  if (videoCam) {
    container.removeChild(videoCam);
  }
}

function handleOnMessage(e) {
  const msg = JSON.parse(e.data);
  const payload = JSON.parse(msg.payload)

  console.log("ws payload => ", payload)

  switch (msg.event) {         
    case JOIN_SIGNAL:      
      let remoteConn = makeRTCPeer();
      myConns[payload.destID] = {
        peerConn: remoteConn,
        candidates: [],
        userConnected: {
          name: payload.userConnected.name,
          myColor: payload.userConnected.myColor,              
        }
      };

      const remoteStream = new MediaStream();

      remoteConn.ontrack = function(evt) {
        evt.streams[0].getTracks().forEach(track => remoteStream.addTrack(track))
      }
      
      streamChannel = remoteConn.createDataChannel("chat");
      streamChannel.onopen = () => console.log("DataChannel opened!");
      streamChannel.onmessage = function(evt) {
        const data = JSON.parse(evt.data);
        addChatMessage(chatBox, data);
        console.log("message from data channel 2: ", evt.data);
      }
      streamChannel.onerror = e => console.error("Channel error 2:", e);

      remoteConn.onicecandidate = evt => {
        if (evt.candidate) {
          const candidate = JSON.stringify(
            evt.candidate.toJSON()
          );

          if (myConns[payload.destID].candidates.length == 0) {
            myConns[payload.destID].candidates = [];
          }
                        
          myConns[payload.destID].candidates.push(candidate);
        }
      }

      myStreamCam.getTracks().forEach(track => remoteConn.addTrack(track, myStreamCam));
      
      remoteConn.createOffer()
        .then(offer => Promise.all([
          Promise.resolve(offer),
          remoteConn.setLocalDescription(offer),
        ]))
        .then((resolver) => {
          const [offer, _] = resolver;

          ws.send(JSON.stringify({
            event: SEND_OFFER_SIGNAL,
            payload: JSON.stringify({
                type: "offer",
                sdp: btoa(JSON.stringify(offer)),                    
                originID: myLocalPeerID,
                destID: payload.destID,
                userConnected: {
                  name: yourname, // my local name
                  myColor: myLocalColor,
                }
              })
          }));
        }).catch(err => console.error("error offer ", err));  

        addCam(payload.destID, payload.userConnected.name, payload.userConnected.myColor, remoteStream, containerEl);
      break;
    case OFFER_SIGNAL:
      let myPeerConn = makeRTCPeer();          
      myStreamCam.getTracks().forEach(track => myPeerConn.addTrack(track, myStreamCam));

      const remoteStream2 = new MediaStream();
      myPeerConn.ondatachannel = function(evt) {
        streamChannel = evt.channel;
        streamChannel.onmessage = (e) => {
          const data = JSON.parse(e.data);
          console.log("message from data channel 1: ", data);
          addChatMessage(chatBox, data);
        };
        streamChannel.onopen = () => console.log("DataChannel opened on callee!");
        streamChannel.onerror = e => console.error("Channel error 1:", e);
      }

      myPeerConn.ontrack = function(evt) {
        evt.streams[0].getTracks().forEach(track => remoteStream2.addTrack(track))
      }

      myPeerConn.onicecandidate = evt => {            
        if (evt.candidate) {
          const candidate = JSON.stringify(
            evt.candidate.toJSON()
          );
          ws.send(JSON.stringify({
            event: SEND_ICE_CANDIDATE_SIGNAL,
            payload: JSON.stringify({
              candidate: candidate,
              originID: payload.destID,
              destID: payload.originID,
            })
          }));
        }
      }

      const myColor = generateRandomColor();
      
      const sdp = JSON.parse(atob(payload.sdp));

      myPeerConn.setRemoteDescription(
        new RTCSessionDescription(sdp)
      ).then(_ => myPeerConn.createAnswer())
      .then(answer => myPeerConn.setLocalDescription(answer))
      .then(_ => {            
        // send answer
        ws.send(JSON.stringify({
          event: SEND_ANSWER_SIGNAL,
          payload: JSON.stringify({
            type: "answer",
            sdp: btoa(JSON.stringify(myPeerConn.localDescription)),
            originID: payload.destID,
            destID: payload.originID,
            userConnected: {
              name: yourname,
              myColor: myLocalColor,
            }
          })
        }));
      })

      myConns[payload.originID] = {
        peerConn: myPeerConn,
        candidates: [],
        userConnected: {
          name: payload.userConnected.name,
          myColor: payload.userConnected.myColor,
        }
      };

      addCam(payload.originID, payload.userConnected.name, payload.userConnected.myColor, remoteStream2, containerEl);          
      break;
    case ANSWER_SIGNAL:
      const answerSdp = JSON.parse(atob(payload.sdp));

      myConns[payload.originID].peerConn.setRemoteDescription(
        new RTCSessionDescription(answerSdp)
      ).then(_ => {
        
        myConns[payload.originID].candidates.forEach(candidate => {

          ws.send(JSON.stringify({
            event: SEND_ICE_CANDIDATE_SIGNAL,
            payload: JSON.stringify({
              candidate: candidate,
              originID: payload.destID,
              destID: payload.originID,
            })
          }));
        });

        myConns[payload.originID].candidates = [];
      }).catch(err => console.error("error set remote description ", err));
      break;
    case JOINED_SIGNAL:
      myLocalPeerID = payload.originID;          
      break;
    case LEAVE_SIGNAL:          
      if (myConns[payload.destID]) {
        myConns[payload.destID].peerConn.close();
        delete myConns[payload.destID];

        removeCam(payload.destID, containerEl);
      }
      break;
    case ICE_CANDIDATE_SIGNAL:          
      const candidate = JSON.parse(payload.candidate);          
      myConns[payload.originID].peerConn.addIceCandidate(
        new RTCIceCandidate(candidate)
      ).then(_ => {
        console.log("add ice candidate");
      }).catch(err => console.error("error add ice candidate ", err));
      break;
  }
}

function makeRTCPeer(urls = []) {
  if (urls.length == 0 ) {
    urls = ['stun:stun1.l.google.com:19302', "stun:stun2.l.google.com:19302"];
  }

  const config = {
    iceServers: [
      {
        urls: urls,
      },
    ],
  };

  return new RTCPeerConnection(config)
}


function generateRoomID() {
  const roomID = Math.random().toString(36).substring(2, 8);
  return roomID;
}

function main() {      
  yourname = prompt("Enter your name", "dian");
  roomID = prompt("Enter room ID", generateRoomID());

  navigator.mediaDevices.getUserMedia({video: true, audio: true})
    .then(stream => {          
      addCam("my-local", yourname, myLocalColor, stream, containerEl);
      myStreamCam = stream;

       // send to socket
      ws.send(JSON.stringify({
        event: NEW_PARTICIPANT_SIGNAL,
        payload: JSON.stringify({
          name: yourname.trim(),
          myColor: myLocalColor,
          roomID: roomID.trim()
        })
      }));
    }).catch(err => {
      console.error("error getUserMedia ", err)
      ws.close();
    });
}


