package model

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

// contant for 3 type actions
const (
	publish     = "publish"
	subscribe   = "subscribe"
	unsubscribe = "unsubscribe"
)

// a server type to store all subscriptions
type Server struct {
	Subscriptions []Subscription
}

// each subscription consists of topic-name & client
type Subscription struct {
	Topic   string
	Clients *[]Client
}

// each client consists of auto-generated ID & connection
type Client struct {
	ID         string
	Connection *websocket.Conn
}

// type for a valid message.
type Message struct {
	Action  string `json:"action"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

func (s *Server) Send(client *Client, message string) {
	client.Connection.WriteMessage(1, []byte(message))
}

func (s *Server) SendWithWait(client *Client, message string, wg *sync.WaitGroup) {
	client.Connection.WriteMessage(1, []byte(message))
	wg.Done()
}

func (s *Server) RemoveClient(client Client) {
	// Read all subs
	for _, sub := range s.Subscriptions {
		// Read all client
		for i := 0; i < len(*sub.Clients); i++ {
			if client.ID == (*sub.Clients)[i].ID {
				// If found, remove client
				if i == len(*sub.Clients)-1 {
					// if it's stored as the last element, crop the array length
					*sub.Clients = (*sub.Clients)[:len(*sub.Clients)-1]
				} else {
					// if it's stored in between elements, overwrite the element and reduce iterator to prevent out-of-bound
					*sub.Clients = append((*sub.Clients)[:i], (*sub.Clients)[i+1:]...)
					i--
				}
			}
		}
	}
}

func (s *Server) ProcessMessage(client Client, messageType int, payload []byte) *Server {
	m := Message{}
	if err := json.Unmarshal(payload, &m); err != nil {
		s.Send(&client, "Server: Invalid payload")
	}

	switch m.Action {
	case publish:
		s.Publish(m.Topic, []byte(m.Message))
		break

	case subscribe:
		s.Subscribe(&client, m.Topic)
		break

	case unsubscribe:
		s.Unsubscribe(&client, m.Topic)
		break

	default:
		s.Send(&client, "Server: Action unrecognized")
		break
	}

	return s
}

func (s *Server) Publish(topic string, message []byte) {
	var clients []Client

	// get list of clients subscribed to topic
	for _, sub := range s.Subscriptions {
		if sub.Topic == topic {
			clients = append(clients, *sub.Clients...)
		}
	}

	// send to clients
	var wg sync.WaitGroup
	for _, client := range clients {
		wg.Add(1)
		s.SendWithWait(&client, string(message), &wg)
	}

	wg.Wait()
}

func (s *Server) Subscribe(client *Client, topic string) {
	exist := false

	// find existing topics
	for _, sub := range s.Subscriptions {
		// if found, add client
		if sub.Topic == topic {
			exist = true
			*sub.Clients = append(*sub.Clients, *client)
		}
	}

	// else, add new topic & add client to that topic
	if !exist {
		newClient := &[]Client{*client}

		newTopic := &Subscription{
			Topic:   topic,
			Clients: newClient,
		}

		s.Subscriptions = append(s.Subscriptions, *newTopic)
	}
}

func (s *Server) Unsubscribe(client *Client, topic string) {
	// Read all topics
	for _, sub := range s.Subscriptions {
		if sub.Topic == topic {
			// Read all topics' client
			for i := 0; i < len(*sub.Clients); i++ {
				if client.ID == (*sub.Clients)[i].ID {
					// If found, remove client
					if i == len(*sub.Clients)-1 {
						// if it's stored as the last element, crop the array length
						*sub.Clients = (*sub.Clients)[:len(*sub.Clients)-1]
					} else {
						// if it's stored in between elements, overwrite the element and reduce iterator to prevent out-of-bound
						*sub.Clients = append((*sub.Clients)[:i], (*sub.Clients)[i+1:]...)
						i--
					}
				}
			}
		}
	}
}
