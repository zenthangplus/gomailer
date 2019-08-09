package main

import (
	"fmt"
	"github.com/zenthangplus/gomailer"
)

func main() {
	client := gomailer.NewClient(
		"smtp.mailtrap.io",
		465,
		"<your-username>",
		"<your-username>",
		gomailer.EncryptionInsecure,
	)
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
