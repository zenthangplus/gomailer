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

type Encryption int

const EncryptionInsecure Encryption = 0
const EncryptionTls Encryption = 1

// Email Client structure
type Client struct {
	host          string
	port          int
	username      string
	password      string
	encryption    Encryption
	defaultSender *Address
}

// Email client constructor
func NewClient(
	host string, port int, username string, password string, encryption Encryption, defaultSender *Address) *Client {
	return &Client{
		host:          host,
		port:          port,
		username:      username,
		password:      password,
		encryption:    encryption,
		defaultSender: defaultSender,
	}
}

// Factory method for sending email
func (c *Client) Send(message MessageInterface) error {
	switch c.encryption {
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
		c.getAddr(), c.getAuth(), message.GetFrom().String(), c.parseAddress(message.GetTo()), []byte(*body),
	)
}

// Send mail by using tls encryption
func (c *Client) sendTls(message MessageInterface) error {
	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for SMTP servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", c.getAddr(), &tls.Config{
		ServerName:         c.host,
		InsecureSkipVerify: true,
	})
	if err != nil {
		return err
	}

	// Create new SMTP client
	client, err := smtp.NewClient(conn, c.host)
	if err != nil {
		return err
	}

	// Authenticate
	if err = client.Auth(c.getAuth()); err != nil {
		return err
	}

	// Set sender
	if err = client.Mail(message.GetFrom().String()); err != nil {
		return err
	}

	// Set recipients
	if err = client.Rcpt(strings.Join(c.parseAddress(message.GetTo()), ",")); err != nil {
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

// Get server address from host and port
func (c *Client) getAddr() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}

// Get authentication config
func (c *Client) getAuth() smtp.Auth {
	return smtp.PlainAuth("", c.username, c.password, c.host)
}

// Build email headers from message
func (c *Client) getHeaders(message MessageInterface) (*string, error) {
	// Init header variable
	header := message.GetHeaders()
	if header == nil {
		header = make(http.Header)
	}

	// Get email sender
	from := message.GetFrom()
	if from == nil && c.defaultSender != nil {
		from = c.defaultSender
	}

	// Check sender is missing or not?
	if from == nil {
		return nil, errors.New("missing sender")
	}

	// Set sender
	header.Set("From", from.String())

	// Set recipients
	header.Set("To", strings.Join(c.parseAddress(message.GetTo()), ","))

	// Set Cc
	if len(message.GetCc()) > 0 {
		header.Set("Cc", strings.Join(c.parseAddress(message.GetCc()), ","))
	}

	// Set Bcc
	if len(message.GetBcc()) > 0 {
		header.Set("Bcc", strings.Join(c.parseAddress(message.GetBcc()), ","))
	}

	// Set Subject
	header.Set("Subject", message.GetSubject())

	// Set other headers
	header.Set("MIME-Version", "1.0")
	header.Set("Content-Type", "text/html; charset=\"utf-8\"")
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

func (c *Client) parseAddress(addressList []*Address) []string {
	recipients := make([]string, 0)
	for _, to := range addressList {
		recipients = append(recipients, to.String())
	}
	return recipients
}
