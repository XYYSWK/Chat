package routers

type routers struct {
	User        user
	Email       email
	Account     account
	Application application
	Message     message
	File        file
	Chat        ws
	Setting     setting
	Group       group
	Notify      notify
}

var Routers = new(routers)
