package websocket

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// constants for action type
const (
	publish     = "publish"
	subscribe   = "subscribe"
	unsubscribe = "unsubscribe"
)

// constants for server message
const (
	errInvalidMessage       = "Server: Invalid msg"
	errActionUnrecognizable = "Server: Action unrecognized"
)

// Server is the struct to handle the Server functions & manage the Subscriptions
type Server struct {
	Subscriptions Subscription
}

// Send simply sends message to the websocket client
func (s *Server) Send(conn *websocket.Conn, message string) {
	// send simple message
	conn.WriteMessage(1, []byte(message))
}

// SendWithWait sends message to the websocket client using wait group, allowing usage with goroutines
func (s *Server) SendWithWait(conn *websocket.Conn, message string, wg *sync.WaitGroup) {
	// send simple message
	conn.WriteMessage(1, []byte(message))

	// set the task as done
	wg.Done()
}

// RemoveClient removes the clients from the server subscription map
func (s *Server) RemoveClient(clientID string) {
	// loop all topics
	for _, client := range s.Subscriptions {
		// delete the client from all the topic's client map
		delete(client, clientID)
	}
}

// ProcessMessage handle message according to the action type
func (s *Server) ProcessMessage(conn *websocket.Conn, clientID string, msg []byte) *Server {
	// parse message
	m := Message{}
	if err := json.Unmarshal(msg, &m); err != nil {
		s.Send(conn, errInvalidMessage)
	}

	// convert all action to lowercase and remove whitespace
	action := strings.TrimSpace(strings.ToLower(m.Action))

	switch action {
	case publish:
		s.Publish(m.Topic, []byte(m.Message))

	case subscribe:
		s.Subscribe(conn, clientID, m.Topic)

	case unsubscribe:
		s.Unsubscribe(clientID, m.Topic)

	default:
		s.Send(conn, errActionUnrecognizable)
	}

	return s
}

// Publish sends a message to all subscribing clients of a topic
func (s *Server) Publish(topic string, message []byte) {
	// if topic does not exist, stop the process
	if _, exist := s.Subscriptions[topic]; !exist {
		return
	}

	// if topic exist
	client := s.Subscriptions[topic]

	// send the message to the clients
	var wg sync.WaitGroup
	for _, conn := range client {
		// add 1 job to wait group
		wg.Add(1)

		// send with goroutines
		go s.SendWithWait(conn, string(message), &wg)
	}

	// wait until all goroutines jobs done
	wg.Wait()
}

// Subscribe adds a client to a topic's client map
func (s *Server) Subscribe(conn *websocket.Conn, clientID string, topic string) {
	// if topic exist, check the client map
	if _, exist := s.Subscriptions[topic]; exist {
		client := s.Subscriptions[topic]

		// if client already subbed, stop the process
		if _, subbed := client[clientID]; subbed {
			return
		}

		// if not subbed, add to client map
		client[clientID] = conn
		return
	}

	// if topic does not exist, create a new topic
	newClient := make(Client)
	s.Subscriptions[topic] = newClient

	// add the client to the topic
	s.Subscriptions[topic][clientID] = conn
}

// Unsubscribe removes a clients from a topic's client map
func (s *Server) Unsubscribe(clientID string, topic string) {
	// if topic exist, check the client map
	if _, exist := s.Subscriptions[topic]; exist {
		client := s.Subscriptions[topic]

		// remove the client from the topic's client map
		delete(client, clientID)
	}
}
