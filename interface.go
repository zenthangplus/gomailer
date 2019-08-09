package gomailer

import (
	"net/http"
)

type MessageInterface interface {
	GetFrom() *Address
	GetTo() []*Address
	GetCc() []*Address
	GetBcc() []*Address
	GetHeaders() http.Header
	GetSubject() string
	GetBody() string
}
