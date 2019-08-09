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
	client := gomailer.NewClient(
		"smtp.mailtrap.io",
		465,
		"<your-username>",
		"<your-password>",
		gomailer.EncryptionInsecure,
	)
	template := gomailer.NewTemplate("template/welcome.html", &gomailer.TemplateConfig{
		LayoutDirectory: "template/layouts",
	})
	templateData := TemplateData{User: &User{
		FirstName: "Go",
		LastName:  "Lang",
	}}
	message, err := gomailer.NewTemplateMessage(&gomailer.Message{
		From:    gomailer.NewAddress("noreply@example.com", "Noreply"),
		To:      []*gomailer.Address{{Address: "user1@example.com"}, {Address: "user2@example.com"}},
		Subject: "Welcome to GoMailer",
	}, template, templateData)
	if err != nil {
		fmt.Print("Err could not make template message: ", err)
	} else {
		if err := client.Send(message); err != nil {
			fmt.Print("Err could not send mail: ", err)
		} else {
			fmt.Print("Success")
		}
	}
}
