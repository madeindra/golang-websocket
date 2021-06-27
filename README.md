# Websocket Server in Go

## Running the server
1. Clone this repository

2. Mount the repository & run this command to install dependencies
```
go get
```

3. Run the websocket server
```
go run main.go
```

4. Websocket server will be running on `localhost:8080`

## Using this server with client
1. After running the server, open your Websocket client. If you don't have any, try `Websocket King` extension for chrome.

2. Connect to `ws://localhost:8080/socket`, you will get an be greeted by the server and.
```
Server: Welcome! Your ID is f0ab664a-5af3-4f8d-8afe-eb93085267e4
```

3. To subscribe to a topic, send this payload (topic can be anything)
```
{
  "action": "subscribe",
  "topic": "world"
}
```

4. To send a message to a specific topic, send payload in this format
```
{
  "action": "publish",
  "topic": "world",
  "message": "Hello world!"
}
```

5. To step receiving message, send this payload (topic can be anything)
```
{
  "action": "unsubscribe",
  "topic": "world"
}
```

## Credit
This repository is a modified version of [Golang-PubSub by @tabvn](https://github.com/tabvn/golang-pubsub-youtube)