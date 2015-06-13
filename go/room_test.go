package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestJoinAddsClientToTheList(t *testing.T) {
	testClients := make(map[*client]bool)
	testClient := getTestClient("John Doe")
	joinRoom(testClient, testClients)
	assert.Len(t, testClients, 1)
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

func getTestClient(name string) *client {
	return &client{
		socket: nil,
		send: make(chan *message),
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