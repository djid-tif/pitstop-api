package schemas

import "html/template"

type Email struct {
	Username    string
	PitStopLink template.URL
	SendingDate string
}

type ConfirmEmail struct {
	Username    string
	ConfirmLink template.URL
	PitStopLink template.URL
	SendingDate string
}

type ResetPassword struct {
	Username    string
	ResetLink   template.URL
	PitStopLink template.URL
	SendingDate string
}

type ReplyToTopic struct {
	Email
	TopicName   string
	ReplierName string
	TopicLink   string
}
