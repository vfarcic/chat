package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestJoinAddsClientToTheList(t *testing.T) {
	room := newRoom()
	testClient := getTestClient("John Doe")
	room.joinRoom(testClient)
	assert.Len(t, room.clients, 1)
}

//func TestJoinShouldSendJoinMessageToAllClients(t *testing.T) {
//	testClients := make(map[*client]bool)
//	client1 := getTestClient("John Doe")
//	client2 := getTestClient("John Doe the Second")
//	testClients[client1] = true
//	testClients[client2] = true
//	newClient := getTestClient("John Doe the Third")
//	go joinRoom(newClient, testClients)
//	expected := make(chan *message)
//	joinMessage := &message{
//		Name: newClient.name,
//		Type: MessageTypeJoin,
//	}
//	assert.Equal(t, expected, newClient.send)
//}

// Helper

func getTestClient(name string) *Client {
	return &Client{
		socket: nil,
		send: make(chan *Message),
		room: nil,
		name: "Viktor Farcic",
		avatarlURL: "http://example.com/avatar",
	}
}

//type MockedClient struct {
//	mock.Mock
//	socket *websocket.Conn
//	send chan *message
//	room *room
//	name string
//	avatarlURL string
//}

//func (m MockedClient) send(user MongoUser) error {
//	ret := m.Called(user)
//	return ret.Error(0)
//}