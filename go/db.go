package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type MongoDB struct {
}

type DB interface {
	Save(msg Message) error
	Drop() error
	GetAll() error
}

func (db MongoDB) Save(msg Message) error {
	session := getSession()
	defer session.Close()
	c := getCollection(session)
	err := c.Insert(msg)
	return err
}

func (db MongoDB) Drop() error {
	session := getSession()
	defer session.Close()
	c := getCollection(session)
	err := c.DropCollection()
	return err
}

func (db MongoDB) GetAll() ([]Message, error) {
	session := getSession()
	defer session.Close()
	c := getCollection(session)
	messages := []Message{}
	err := c.Find(bson.M{}).All(&messages)
	return messages, err
}

func getSession() *mgo.Session {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	return session
}

func getCollection(session *mgo.Session) *mgo.Collection {
	return session.DB("chat").C("messages")
}
