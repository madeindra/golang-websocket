package model

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

// contant for 3 type event
const (
	publish     = "publish"
	subscribe   = "subscribe"
	unsubscribe = "unsubscribe"
)

// all clients & subscriptions are stored in a variable of this type
type PubSub struct {
	Clients       []Client
	Subscriptions []subscription
}

// each client consists of auto-generated ID & connection
type Client struct {
	ID         string
	Connection *websocket.Conn
}

// each subscription consists of topic-name & client
type subscription struct {
	Topic  string
	client *Client
}

// type for a valid message.
type Message struct {
	Action  string `json:"action"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

// called when a new client connected to the socket
func (ps *PubSub) AddClient(client Client) *PubSub {
	ps.Clients = append(ps.Clients, client)
	return ps
}

// called when a client disconnected from the socket
func (ps *PubSub) RemoveClient(client Client) *PubSub {
	// remove all subscriptions by this client
	for i := 0; i < len(ps.Subscriptions); i++ {
		sub := ps.Subscriptions[i]
		if client.ID == sub.client.ID {
			if i == len(ps.Subscriptions)-1 {
				// if it's stored as the last element, crop the array length
				ps.Subscriptions = ps.Subscriptions[:len(ps.Subscriptions)-1]
			} else {
				// if it's stored in between elements, overwrite the element and reduce iterator to prevent out-of-bound
				ps.Subscriptions = append(ps.Subscriptions[:i], ps.Subscriptions[i+1:]...)
				i--
			}
		}
	}

	// remove this client from the clients list
	for i := 0; i < len(ps.Clients); i++ {
		c := ps.Clients[i]
		if c.ID == client.ID {
			if i == len(ps.Clients)-1 {
				// if it's stored as the last element, crop the array length by 1
				ps.Clients = ps.Clients[:len(ps.Clients)-1]
			} else {
				// if it's stored in between elements, overwrite it by the next element and reduce iterator to prevent out-of-bound
				ps.Clients = append(ps.Clients[:i], ps.Clients[i+1:]...)
				i--
			}
		}
	}

	return ps
}

// called when 'publish' action is called, sends a message to all topic subscriber
func (ps *PubSub) Publish(topic string, message []byte) {
	subscriptions := ps.GetSubscriptions(topic, nil)

	for _, sub := range subscriptions {
		sub.client.Send(message)
	}
}

// called when needed e.g. wrong action is called, sends a message back to the action caller
func (ps *PubSub) BounceBack(client *Client, message string) {
	client.Send([]byte(message))
}

// called when 'subscribe' action is called, adds a new topic subscription
func (ps *PubSub) Subscribe(client *Client, topic string) *PubSub {
	clientSubs := ps.GetSubscriptions(topic, client)

	if len(clientSubs) > 0 {
		// client has subscribed this topic before
		return ps
	}

	newSubscription := subscription{
		Topic:  topic,
		client: client,
	}

	ps.Subscriptions = append(ps.Subscriptions, newSubscription)

	return ps
}

// called when 'unsubscribe' action is called, remove a topic subscription
func (ps *PubSub) Unsubscribe(client *Client, topic string) *PubSub {
	// checks if client subscribed to the topic
	for i := 0; i < len(ps.Subscriptions); i++ {
		sub := ps.Subscriptions[i]
		if sub.client.ID == client.ID && sub.Topic == topic {
			// found this subscription from client, we need remove it
			if i == len(ps.Subscriptions)-1 {
				ps.Subscriptions = ps.Subscriptions[:len(ps.Subscriptions)-1]
			} else {
				ps.Subscriptions = append(ps.Subscriptions[:i], ps.Subscriptions[i+1:]...)
				i--
			}
		}
	}

	return ps
}

// used by publish & bounceback, this is a basic websocket message sending function
func (client *Client) Send(message []byte) error {
	return client.Connection.WriteMessage(1, message)
}

// used by publish & subscribe, this either add a new topic or get a list of subscription
func (ps *PubSub) GetSubscriptions(topic string, client *Client) []subscription {
	var subscriptionList []subscription

	for _, subscription := range ps.Subscriptions {
		if client != nil {
			// if no client is provided, then give all subscription
			if subscription.client.ID == client.ID && subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)
			}
		} else {
			// else, then give client's subscription
			if subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)
			}
		}
	}

	return subscriptionList
}

// used for message handling, runs function(s) according to actions sent
func (ps *PubSub) ProcessMessage(client Client, messageType int, payload []byte) *PubSub {
	m := Message{}
	if err := json.Unmarshal(payload, &m); err != nil {
		ps.BounceBack(&client, "Server: Invalid payload")
	}

	switch m.Action {
	case publish:
		ps.Publish(m.Topic, []byte(m.Message))
		break

	case subscribe:
		ps.Subscribe(&client, m.Topic)
		break

	case unsubscribe:
		ps.Unsubscribe(&client, m.Topic)
		break

	default:
		ps.BounceBack(&client, "Server: Action unrecognized")
		break
	}

	return ps
}
