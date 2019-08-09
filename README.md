# Go Mailer

Mailer package for Golang supports sending SMTP email using [template](https://golang.org/pkg/html/template/).

## Installation

Run the following command to install this package:

```
$ go get -u github.com/zenthangplus/gomailer
``` 

## Usage

The following example will help you known how to use this package:

```go
package main

import (
	"fmt"
	"github.com/zenthangplus/gomailer"
)

type User struct {
	FirstName string
	LastName  string
}

type TemplateData struct {
	User *User
}

func main() {
	// Create the email client
	client := gomailer.NewClient(
		"smtp.example.com", // SMTP Host
		465,// SMTP Port
		"<your-username>", // SMTP Username
		"<your-password>", // SMTP Password
		gomailer.EncryptionTls, // SMTP Encryption
	)
	
	// Create the template
	template := gomailer.NewTemplate("template/welcome.html", &gomailer.TemplateConfig{
		LayoutDirectory: "template/layouts",
	})
	
	// Create the template data
	templateData := TemplateData{User: &User{
		FirstName: "Go",
		LastName:  "Lang",
	}}
	
	// Create the message using the template
	message, err := gomailer.NewTemplateMessage(&gomailer.Message{
		From:    gomailer.NewAddress("noreply@example.com", "Example"),
		To:      []*gomailer.Address{{Address: "user1@example.com"}, {Address: "user2@example.com"}},
		Subject: "Welcome to GoMailer",
	}, template, templateData)
	if err != nil {
		fmt.Print("Err could not make template message: ", err)
		return
	}
	
	// Try to send the email
	if err := client.Send(message); err != nil {
        fmt.Print("Err could not send mail: ", err)
    } else {
        fmt.Print("Success")
    }
}

```

> Go to [example](/example) directory to see more examples.