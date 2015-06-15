package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
)

type Room struct {
	forward chan *Message
	join chan *Client
	leave chan *Client
	clients map[*Client]bool
}

func newRoom() *Room {
	return &Room{
		forward: make(chan *Message),
		join: make(chan *Client),
		leave: make(chan *Client),
		clients: make(map[*Client]bool),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.join:
			r.joinRoom(client)
		case client := <-r.leave:
			r.leaveRoom(client)
		case msg := <-r.forward:
			r.sendMessage(msg, MongoDB{})
		}
	}
}

func (r *Room) joinRoom(joinClient *Client) {
	if !r.clients[joinClient] {
		for clientToSend := range r.clients {
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
	r.clients[joinClient] = true
	// TODO: Send all message from history
//	messages := db
//	clientToSend.send <- msg
}

func (r *Room) leaveRoom(leaveClient *Client) {
	if r.clients[leaveClient] {
		for clientToSend := range r.clients {
			msg := &Message{
				Name: leaveClient.name,
				Type: MessageTypeLeave,
			}
			clientToSend.send <-msg
		}
		log.Println(leaveClient.name, "left")
	}
	delete(r.clients, leaveClient)
	close(leaveClient.send)
}

func (r *Room) sendMessage(msg *Message, db DB) error {
	log.Println("Messsage received from", msg.Name, ":\n", msg.Message)
	for clientToSend := range r.clients {
		select {
		case clientToSend.send <- msg:
		// Send the message
		default:
			delete(r.clients, clientToSend)
			close(clientToSend.send)
			log.Println(" -- failed to send")
		}
	}
	err := db.Save(*msg)
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

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
