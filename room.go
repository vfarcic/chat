package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
)

type room struct {
	forward chan *message
	join chan *client
	leave chan *client
	clients map[*client]bool
}

func newRoom() *room {
	return &room {
		forward: make(chan *message),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			if !r.clients[client] {
				for clientToSend := range r.clients {
					msg := &message{
						Name: client.name,
						Type: MessageTypeJoin,
					}
					clientToSend.send <-msg
				}
				log.Println(client.name, "joined")
			}
			r.clients[client] = true
		case client := <-r.leave:
			if r.clients[client] {
				for clientToSend := range r.clients {
					msg := &message{
						Name: client.name,
						Type: MessageTypeLeave,
					}
					clientToSend.send <-msg
				}
				log.Println(client.name, "left")
			}
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			log.Println("Messsage received from", msg.Name, ":\n", msg.Message)
			for client := range r.clients {
				select {
				case client.send <-msg:
					// Send the message
				default:
					delete(r.clients, client)
					close(client.send)
					log.Println(" -- failed to send")
				}
			}
		}
	}
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
	client := &client {
		socket: socket,
		send: make(chan *message, messageBufferSize),
		room: r,
		name: authName,
		avatarlURL: authAvatarURL,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
