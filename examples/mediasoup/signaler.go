package mediasoup

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/ice"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Signaler struct {
	url        string
	conn       *websocket.Conn
	done       chan struct{}
	idMapMutex *sync.Mutex
	idMap      map[int]chan []byte
	rtcPeer    *webrtc.RTCPeerConnection
	rand       *rand.Rand
}

type request struct {
	Request bool        `json:"request"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Data    interface{} `json:"data"`
}

type response struct {
	Response bool            `json:"response"`
	ID       int             `json:"id"`
	Ok       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
}

func NewSignaler(rawUrl string) *Signaler {

	rs := rand.NewSource(time.Now().UnixNano())

	s := &Signaler{}
	s.idMapMutex = &sync.Mutex{}
	s.rand = rand.New(rs)
	s.done = make(chan struct{})
	s.idMap = make(map[int]chan []byte)
	s.url = rawUrl
	s.rtcPeer = &webrtc.RTCPeerConnection{}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	s.rtcPeer.OnICEConnectionStateChange = func(connectionState ice.ConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	}

	return s
}

func (s *Signaler) Connect() error {

	header := make(http.Header)
	header.Add("Sec-WebSocket-Protocol", "protoo")

	c, _, err := websocket.DefaultDialer.Dial(s.url, header)
	if err != nil {
		return errors.Wrap(err, "Failed to connect to signaler")
	}

	s.conn = c

	go func() {
		defer close(s.done)
		for {
			_, message, err := c.ReadMessage()
			t := &response{}
			err = json.Unmarshal(message, t)
			if err != nil {
				log.Println("read:", err)
				return
			}

			if v, ok := s.idMap[t.ID]; ok {
				v <- message
			} else {
				log.Printf("unknown ID %v, recv: %s", t.ID, message)
			}
		}
	}()

	return nil
}

func (s *Signaler) Request(data interface{}) (chan []byte, error) {

	var newId int
	{
		s.idMapMutex.Lock()
		defer s.idMapMutex.Unlock()
		// Create a random, unused id
		ok := true
		for ok {
			newId = int(s.rand.Int31())
			_, ok = s.idMap[newId]
		}

		s.idMap[newId] = make(chan []byte)
	}

	r := request{
		Request: true,
		Method:  "mediasoup-request",
		ID:      newId,
		Data:    data,
	}

	fmt.Printf("Sending request: %v\n", r)

	err := s.conn.WriteJSON(r)

	if err != nil {
		delete(s.idMap, newId)
		return nil, errors.Wrap(err, "Unable to send request")
	}

	return s.idMap[newId], nil
}

func (s *Signaler) Close() {
	close(s.done)
}
