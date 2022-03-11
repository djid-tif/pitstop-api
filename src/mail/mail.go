package mail

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
	"path"
)

// SendMail sends a mail to the given address with the given subject and body.
// The name parameter is the name of the receiver and the bodyType parameter is the content-type of the body
func SendMail(name, address, subject, body, bodyType string) error {
	dialer := gomail.NewDialer(credentials.Host, credentials.Port, credentials.User, credentials.Password)

	message := gomail.NewMessage()
	message.SetAddressHeader("From", "no-reply@pitstop.com", "PitStop")
	message.SetAddressHeader("To", address, name)
	message.SetHeader("Subject", "PitStop - "+subject)

	message.SetBody(bodyType, body)

	return dialer.DialAndSend(message)
}

// SendTemplate sends a mail to the given address with the given subject.
// The mail's body is the given template executed with the given data.
func SendTemplate(name, address, subject, templateName string, data interface{}) error {
	templateFile, err := templates.ReadFile(path.Join("templates", templateName+".tmpl"))
	if err != nil {
		return fmt.Errorf("échec du chargement de la template '%s': %v\n", templateName, err)
	}

	tmpl, err := template.New(templateName).Parse(string(templateFile))
	if err != nil {
		return fmt.Errorf("échec du parsing de la template '%s': %v\n", templateName, err)
	}

	receiver := writeReceiver{buffer: []byte{}}
	err = tmpl.Execute(&receiver, data)
	if err != nil {
		return fmt.Errorf("échec de l'exécution de la template '%s': %v\n", templateName, err)
	}

	return SendMail(name, address, subject, receiver.Export(), "text/html")
}
