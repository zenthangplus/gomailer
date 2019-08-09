package gomailer

import (
	"net/http"
	"net/mail"
)

type Address = mail.Address

func NewAddress(address string, name string) *Address {
	return &Address{
		Address: address,
		Name:    name,
	}
}

type Message struct {
	From    *Address
	To      []*Address
	Cc      []*Address
	Bcc     []*Address
	Headers http.Header
	Subject string
	Body    string
}

func NewMessage(from *Address, to []*Address, subject string, body string) *Message {
	return &Message{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    body,
	}
}

func (e *Message) GetFrom() *Address {
	return e.From
}

func (e *Message) GetTo() []*Address {
	return e.To
}

func (e *Message) GetCc() []*Address {
	return e.Cc
}

func (e *Message) GetBcc() []*Address {
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
