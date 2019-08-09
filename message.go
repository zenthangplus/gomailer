package gomailer

import (
	"net/http"
	"net/mail"
)

type Message struct {
	From    *mail.Address
	To      []string
	Cc      []string
	Bcc     []string
	Headers http.Header
	Subject string
	Body    string
}

func (e *Message) GetFrom() *mail.Address {
	return e.From
}

func (e *Message) GetTo() []string {
	return e.To
}

func (e *Message) GetCc() []string {
	return e.Cc
}

func (e *Message) GetBcc() []string {
	return e.Bcc
}

func (e *Message) GetHeaders() http.Header {
	return e.Headers
}

func (e *Message) GetSubject() string {
	return e.Subject
}

func (e *Message) GetBody() string {
	return e.Body
}
