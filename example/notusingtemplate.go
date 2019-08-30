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

	// Try to send the email with a simple message
	if err := client.Send(&gomailer.Message{
		From:    gomailer.NewAddress("noreply@example.com", "Noreply"),
		To:      []*gomailer.Address{{Address: "user1@example.com"}, {Address: "user2@example.com"}},
		Subject: "Welcome to GoMailer",
		Body:    "You are using GoMailer",
	}); err != nil {
		fmt.Print("Err could not send mail: ", err)
	} else {
		fmt.Print("Success")
	}
}
