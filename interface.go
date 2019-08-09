package gomailer

import (
	"net/http"
	"net/mail"
)

type MessageInterface interface {
	GetFrom() *mail.Address
	GetTo() []string
	GetCc() []string
	GetBcc() []string
	GetHeaders() http.Header
	GetSubject() string
	GetBody() string
}
