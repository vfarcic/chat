package main

type message struct {
	Name string
	Message string
	Date string
	Time string
	AvatarURL string
	Type string
}

const MessageTypeMessage = "message"
const MessageTypeJoin = "join"
const MessageTypeLeave = "leave"
