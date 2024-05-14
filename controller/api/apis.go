package api

import "Chat/controller/api/chat"

type apis struct {
	User        user
	Email       email
	Account     account
	Application application
	File        file
	Message     message
	Chat        chat.Group
	Setting     setting
	Group       group
	Notify      notify
}

var Apis = new(apis)
