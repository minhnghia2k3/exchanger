package mail

import (
	"bytes"
	"embed"
	gomail "gopkg.in/mail.v2"
	"html/template"
	"time"
)

//go:embed templates
var templateFS embed.FS

type Mailer struct {
	Dialer *gomail.Dialer
	Sender string
}

func NewMailer(sender string, host string, port int, username, password string) *Mailer {
	dialer := gomail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return &Mailer{
		Dialer: dialer,
		Sender: sender,
	}
}

func (m *Mailer) Send(recipient, templateFile string, data any) error {
	// Parse template
	tmpl, err := template.ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// Execute subject
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// Execute plain body
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	// Execute html body
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	message := gomail.NewMessage()

	message.SetHeader("From", m.Sender)
	message.SetHeader("To", recipient)
	message.SetHeader("Subject", subject.String())

	message.SetBody("text/html", htmlBody.String())
	message.AddAlternative("text/plain", plainBody.String())

	retry := 3
	for i := 0; i < retry; i++ {
		err = m.Dialer.DialAndSend(message)
		if err == nil {
			return nil
		}

		// If there is error, resend email
		time.Sleep(500 * time.Millisecond)
	}

	return err
}
