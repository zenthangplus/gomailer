# Go Mailer

Mailer package for Golang supports sending SMTP email using [html/template](https://golang.org/pkg/html/template/).

## Installation

Run the following command to install this package:

```
$ go get -u github.com/zenthangplus/gomailer
``` 

## Usage

The following example will help you know how to use this package:

> Note: Go to [example](/example) directory to see more examples.

```go
package main

import (
    "fmt"
    "github.com/zenthangplus/gomailer"
)

func main() {
    // 1. Create the email client
    client := gomailer.NewClient(
        "smtp.example.com",     // Host
        465,                    // Port
        "<your-username>",      // Username
        "<your-password>",      // Password
        gomailer.EncryptionTls, // Encryption
    )
	
    // 2. Create the template
    // The first argument is the template path
    // The second argument is the config for your template. 
    // Such as: LayoutFiles, LayoutDirectory, LayoutExtension (available when using LayoutDirectory).
    // These layout files will be included to your template.
    // @see https://golang.org/pkg/html/template/
    template := gomailer.NewTemplate("template/welcome.html", &gomailer.TemplateConfig{
        LayoutDirectory: "template/layouts",
    })
	
    // 3. Create the template data
    // You can use Map or Struct here. It is flexible.
    templateData := map[string]interface{}{
        "User": map[string]string{
            "FirstName": "Go",
            "LastName":  "Lang",
        },
    }
	
    // 4. Create the message using the template
    message, err := gomailer.NewTemplateMessage(&gomailer.Message{
        From:    gomailer.NewAddress("noreply@example.com", "Example"),
        To:      []*gomailer.Address{{Address: "user1@example.com"}, {Address: "user2@example.com"}},
        Subject: "Welcome to GoMailer",
    }, template, templateData)
    if err != nil {
        fmt.Print("Err could not make template message: ", err)
        return
    }
	
    // 5. Try to send the message
    if err := client.Send(message); err != nil {
        fmt.Print("Err could not send mail: ", err)
    } else {
        fmt.Print("Success")
    }
}
```
