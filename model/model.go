package model

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

const (
	publish     = "publish"
	subscribe   = "subscribe"
	unsubscribe = "unsubscribe"
)

type PubSub struct {
	Clients       []Client
	Subscriptions []subscription
}

type Client struct {
	ID         string
	Connection *websocket.Conn
}

type subscription struct {
	Topic  string
	client *Client
}

type Message struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type Chats struct {
	TargetID string `json:"targetId"`
	Message  string `json:"message"`
}

func (ps *PubSub) AddClient(client Client) *PubSub {
	ps.Clients = append(ps.Clients, client)
	return ps
}

func (ps *PubSub) RemoveClient(client Client) *PubSub {
	// remove all subscriptions by this client
	for i := 0; i < len(ps.Subscriptions); i++ {
		sub := ps.Subscriptions[i]
		if client.ID == sub.client.ID {
			if i == len(ps.Subscriptions)-1 {
				ps.Subscriptions = ps.Subscriptions[:len(ps.Subscriptions)-1]
			} else {
				ps.Subscriptions = append(ps.Subscriptions[:i], ps.Subscriptions[i+1:]...)
				i--
			}
		}
	}

	// remove this client from the list
	for i := 0; i < len(ps.Clients); i++ {
		c := ps.Clients[i]
		if c.ID == client.ID {
			if i == len(ps.Clients)-1 {
				ps.Clients = ps.Clients[:len(ps.Clients)-1]
			} else {
				ps.Clients = append(ps.Clients[:i], ps.Clients[i+1:]...)
				i--
			}
		}
	}

	return ps
}

func (ps *PubSub) Publish(topic string, message []byte) {
	subscriptions := ps.GetSubscriptions(topic, nil)

	for _, sub := range subscriptions {
		sub.client.Send(message)
	}
}

func (ps *PubSub) BounceBack(client *Client, message string) {
	client.Send([]byte(message))
}

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

func (ps *PubSub) Unsubscribe(client *Client, topic string) *PubSub {
	for i := 0; i < len(ps.Subscriptions); i++ {
		sub := ps.Subscriptions[i]
		if sub.client.ID == client.ID && sub.Topic == topic {
			// found this subscription from client and we do need remove it
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

func (client *Client) Send(message []byte) error {
	return client.Connection.WriteMessage(1, message)
}

func (ps *PubSub) GetSubscriptions(topic string, client *Client) []subscription {
	var subscriptionList []subscription

	for _, subscription := range ps.Subscriptions {
		if client != nil {
			if subscription.client.ID == client.ID && subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)
			}
		} else {
			if subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)
			}
		}
	}

	return subscriptionList
}

func (ps *PubSub) ProcessMessage(client Client, messageType int, payload []byte) *PubSub {
	m := Message{}
	if err := json.Unmarshal(payload, &m); err != nil {
		ps.BounceBack(&client, "Server: Failed binding action")
	}

	switch m.Action {
	case publish:
		ch := Chats{}
		if err := json.Unmarshal(m.Data, &ch); err != nil {
			ps.BounceBack(&client, "Server: Failed binding data")
			break
		}

		ps.Publish(ch.TargetID, []byte(ch.Message))
		break

	case subscribe:
		ps.BounceBack(&client, "Server: Welcome! Your ID is "+client.ID)
		ps.Subscribe(&client, client.ID)
		break

	case unsubscribe:
		ps.Unsubscribe(&client, client.ID)
		break

	default:
		ps.BounceBack(&client, "Server: Action unrecognized")
		break
	}

	return ps
}
