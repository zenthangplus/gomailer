package gomailer

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
)

// List of auth types
const (
	AuthPlain   AuthType = 0
	AuthCramMd5 AuthType = 1
)

// List of encryption types
const (
	EncryptionInsecure string = ""
	EncryptionTls      string = "tls"
)

// Auth type
type AuthType int

// Encryption type
type Encryption int

// Email Client structure
type Client struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Encryption string
	AuthType   AuthType
	addr       string
	auth       smtp.Auth
}

// Factory method for sending email
func (c *Client) Send(message MessageInterface) error {
	// Get server address from host and port
	c.addr = fmt.Sprintf("%s:%d", c.Host, c.Port)

	// Get SMTP Auth config
	switch c.AuthType {
	case AuthPlain:
		c.auth = smtp.PlainAuth("", c.Username, c.Password, c.Host)
	case AuthCramMd5:
		c.auth = smtp.CRAMMD5Auth(c.Username, c.Password)
	}

	// Send the message
	switch c.Encryption {
	case EncryptionTls:
		return c.sendTls(message)
	case EncryptionInsecure:
		return c.sendInsecure(message)
	default:
		return errors.New("invalid encryption")
	}
}

// Send mail by using non-encryption
func (c *Client) sendInsecure(message MessageInterface) error {
	// Build the message body
	body, err := c.getBody(message)
	if err != nil {
		return err
	}
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	return smtp.SendMail(
		c.addr,
		c.auth,
		message.GetFrom().Address,
		c.parseAddressList(message.GetTo(), false),
		[]byte(*body),
	)
}

// Send mail by using tls encryption
func (c *Client) sendTls(message MessageInterface) error {
	// Connect to SMTP server
	client, err := smtp.Dial(c.addr)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Close()
	}()

	if err := client.StartTLS(&tls.Config{
		ServerName:         c.Host,
		InsecureSkipVerify: true,
	}); err != nil {
		return err
	}

	// Authenticate
	if err = client.Auth(c.auth); err != nil {
		return err
	}

	// Set sender
	if err = client.Mail(message.GetFrom().Address); err != nil {
		return err
	}

	// Set recipients
	if err = client.Rcpt(strings.Join(c.parseAddressList(message.GetTo(), false), ", ")); err != nil {
		return err
	}

	// Get body from message
	body, err := c.getBody(message)
	if err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write([]byte(*body)); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	return client.Quit()
}

// Build email headers from message
func (c *Client) getHeaders(message MessageInterface) (*string, error) {
	// Init header variable
	header := message.GetHeaders()
	if header == nil {
		header = make(http.Header)
	}

	// Check sender is missing or not?
	if message.GetFrom() == nil {
		return nil, errors.New("missing sender")
	}

	// Set sender
	header.Set("From", message.GetFrom().String())

	// Set recipients
	header.Set("To", strings.Join(c.parseAddressList(message.GetTo(), true), ", "))

	// Set Cc
	if len(message.GetCc()) > 0 {
		header.Set("Cc", strings.Join(c.parseAddressList(message.GetCc(), true), ", "))
	}

	// Set Bcc
	if len(message.GetBcc()) > 0 {
		header.Set("Bcc", strings.Join(c.parseAddressList(message.GetBcc(), true), ", "))
	}

	// Set Subject
	header.Set("Subject", message.GetSubject())

	// Set default mime version
	if _, exists := header["MIME-Version"]; exists == false {
		header.Set("MIME-Version", "1.0")
	}

	// Set default content type
	if _, exists := header["Content-Type"]; exists == false {
		header.Set("Content-Type", "text/html; charset=\"utf-8\"")
	}

	// Set encoding type for this mail
	header.Set("Content-Transfer-Encoding", "base64")

	// Generate headers string
	var headers bytes.Buffer
	if err := header.Write(&headers); err != nil {
		return nil, err
	}
	headerStr := headers.String()
	return &headerStr, nil
}

// Build message body
func (c *Client) getBody(message MessageInterface) (*string, error) {
	// Build message's headers
	headers, err := c.getHeaders(message)
	if err != nil {
		return nil, err
	}

	// Make raw message
	body := *headers + "\r\n" + base64.StdEncoding.EncodeToString([]byte(message.GetBody()))
	return &body, nil
}

// Parse email address list to list of address (in string)
func (c *Client) parseAddressList(addressList []*Address, formatRfc5322 bool) []string {
	recipients := make([]string, 0)
	for _, to := range addressList {
		formattedAddr := to.Address
		if formatRfc5322 == true {
			formattedAddr = to.String()
		}
		recipients = append(recipients, formattedAddr)
	}
	return recipients
}
