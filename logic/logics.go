package logic

type logics struct {
	User        user
	Email       email
	Auto        auto
	Account     account
	Application application
	File        file
	Message     message
	Setting     setting
	Group       group
	Notify      notify
}

var Logics = new(logics)
