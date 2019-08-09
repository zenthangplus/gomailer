package gomailer

import (
	"net/http"
)

type MessageInterface interface {
	GetFrom() *Address
	GetTo() []string
	GetCc() []string
	GetBcc() []string
	GetHeaders() http.Header
	GetSubject() string
	GetBody() string
}
