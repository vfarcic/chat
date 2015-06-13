package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
)

type room struct {
	forward chan *Message
	join chan *Client
	leave chan *Client
	clients map[*Client]bool
}

func newRoom() *room {
	return &room {
		forward: make(chan *Message),
		join: make(chan *Client),
		leave: make(chan *Client),
		clients: make(map[*Client]bool),
	}
}

func (r *room) Run() {
	for {
		select {
		case client := <-r.join:
			joinRoom(client, r.clients)
		case client := <-r.leave:
			leaveRoom(client, r.clients)
		case msg := <-r.forward:
			sendMessage(msg, r.clients, MongoDB{})
		}
	}
}

func joinRoom(joinClient *Client, clients map[*Client]bool) {
	if !clients[joinClient] {
		for clientToSend := range clients {
			msgNewClient := &Message{
				Name: joinClient.name,
				Type: MessageTypeJoin,
			}
			msgOldClient := &Message{
				Name: clientToSend.name,
				Type: MessageTypeJoin,
			}
			clientToSend.send <-msgNewClient
			joinClient.send <-msgOldClient
		}
		log.Println(joinClient.name, "joined")
	}
	clients[joinClient] = true
}

func leaveRoom(leaveClient *Client, clients map[*Client]bool) {
	if clients[leaveClient] {
		for clientToSend := range clients {
			msg := &Message{
				Name: leaveClient.name,
				Type: MessageTypeLeave,
			}
			clientToSend.send <-msg
		}
		log.Println(leaveClient.name, "left")
	}
	delete(clients, leaveClient)
	close(leaveClient.send)
}

func sendMessage(msg Message, clients map[*Client]bool, db DB) error {
	log.Println("Messsage received from", msg.Name, ":\n", msg.Message)
	for clientToSend := range clients {
		select {
		case clientToSend.send <- msg:
		// Send the message
		default:
			delete(clients, clientToSend)
			close(clientToSend.send)
			log.Println(" -- failed to send")
		}
	}
	err := db.Save(msg)
	return err
}

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)
var upgrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
	ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	authNameCookie, err := req.Cookie("authName")
	var authName = "John Doe"
	if err == nil {
		authName = authNameCookie.Value
	}
	authAvatarURLCookie, err := req.Cookie("authAvatarURL")
	authAvatarURL := ""
	if err == nil {
		authAvatarURL = authAvatarURLCookie.Value
	}
	client := &Client{
		socket: socket,
		send: make(chan *Message, messageBufferSize),
		room: r,
		name: authName,
		avatarlURL: authAvatarURL,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
