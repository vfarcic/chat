package main

import (
	"github.com/gorilla/websocket"
	"time"
)

type client struct {
	socket *websocket.Conn
	send chan *message
	room *room
	name string
	avatarlURL string
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.Date = time.Now().Format("Monday, 02-Jan-06")
			msg.Time = time.Now().Format("15:04")
			msg.Name = c.name
			msg.AvatarURL = c.avatarlURL
			msg.Type = MessageTypeMessage
			c.room.forward <- msg
		} else {
			break;
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break;
		}
	}
	c.socket.Close()
}
