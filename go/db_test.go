package main

import (
	"testing"
	"time"
	"labix.org/v2/mgo/bson"
	"github.com/stretchr/testify/assert"
	"fmt"
)

var db = MongoDB{}
var msg = Message{
	Name: "John Doe",
	Message: "Hello World!",
	Date: time.Now().Format("Monday, 02-Jan-06"),
	Time: time.Now().Format("15:04"),
	AvatarURL: "http://example.com/avatar",
	Type: MessageTypeMessage,
}

func TestSaveToMongoDBShouldInsertData(t *testing.T) {
	db.Drop()
	db.Save(msg)

	session := getSession()
	defer session.Close()
	c := getCollection(session)
	actual := Message{}
	err := c.Find(bson.M{"name": msg.Name, "message": msg.Message}).One(&actual)

	assert.Nil(t, err)
	assert.Equal(t, msg, actual)
}

// Drop

func TestDropShouldRemoveCollection(t *testing.T) {
	db.Drop()
	db.Save(msg)

	err := db.Drop()

	session := getSession()
	defer session.Close()
	c := getCollection(session)
	actual := []Message{}
	c.Find(nil).All(&actual)
	assert.Nil(t, err)
	assert.Len(t, actual, 0)
}

// GetAll

func TestGetAllShouldReturnAllRecords(t *testing.T) {
	db.Drop()
	for index := 1; index <= 5; index++ {
		msg.Message = fmt.Sprintln("Message ", index)
		db.Save(msg)
	}

	messages, err := db.GetAll()

	assert.Nil(t, err)
	assert.Len(t, messages, 5)
}
