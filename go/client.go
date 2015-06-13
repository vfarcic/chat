package main

import (
	"github.com/gorilla/websocket"
	"time"
)

type Client struct {
	socket *websocket.Conn
	send chan *Message
	room *room
	name string
	avatarlURL string
}

func (c *Client) read() {
	for {
		var msg *Message
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

func (c *Client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break;
		}
	}
	c.socket.Close()
}
