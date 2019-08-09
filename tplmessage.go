package gomailer

type TemplateMessage struct {
	Message
}

func NewTemplateMessage(message *Message, template *Template, data interface{}) (*TemplateMessage, error) {
	body, err := template.Parse(data)
	if err != nil {
		return nil, err
	}
	message.Body = body
	return &TemplateMessage{
		Message: *message,
	}, nil
}
