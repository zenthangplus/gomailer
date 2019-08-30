package main

import (
	"fmt"
	"github.com/zenthangplus/gomailer"
)

func main() {
	// Create the email client
	client := gomailer.Client{
		Host:       "smtp.example.com",
		Port:       465,
		Username:   "<your-username>",
		Password:   "<your-password>",
		Encryption: gomailer.EncryptionTls,
	}

	// Create the template
	template := gomailer.NewTemplate("template/welcome.html", &gomailer.TemplateConfig{
		LayoutDirectory: "template/layouts",
	})

	// Create the template data
	type User struct {
		FirstName string
		LastName  string
	}
	templateData := map[string]interface{}{
		"User": User{
			FirstName: "Go",
			LastName:  "Lang",
		},
	}

	// Create the message using the template
	message, err := gomailer.NewTemplateMessage(&gomailer.Message{
		From:    gomailer.NewAddress("noreply@example.com", "Noreply"),
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
