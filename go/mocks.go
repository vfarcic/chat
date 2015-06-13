package main
import "github.com/stretchr/testify/mock"

type MockedRoom struct {
	mock.Mock
	forward chan *Message
	join chan *Client
	leave chan *Client
	clients map[*Client]bool
}

func (m *MockedRoom) Run() {
}

func (m *MockedRoom) joinRoom(joinClient *Client) {
}

func (m *MockedRoom) leaveRoom(leaveClient *Client) {
}

func (m *MockedRoom) sendMessage(msg Message, db DB) error {
	ret := m.Called(msg, db)
	return ret.Error(0)
}