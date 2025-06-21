package mailer

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGrid(fromEmail, apiKey string) *SendGridMailer {

	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (mailer *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {

	fmt.Println("0")

	from := mail.NewEmail(FromName, mailer.fromEmail)
	to := mail.NewEmail(username, email)

	fmt.Println(from, to)

	// template parsing and building
	// tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("2")

	// subject := new(bytes.Buffer)
	// err = tmpl.ExecuteTemplate(subject, "subject", data)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("3")

	// body := new(bytes.Buffer)
	// err = tmpl.ExecuteTemplate(subject, "body", data)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("4")

	// message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())
	// message.SetMailSettings(&mail.MailSettings{
	// 	SandboxMode: &mail.Setting{
	// 		Enable: &isSandbox,
	// 	},
	// })

	// fmt.Println("5")

	// for i := 0; i < maxRetries; i++ {
	// 	response, err := mailer.client.Send(message)
	// 	if err != nil {
	// 		log.Printf("Failed to send email to %v, attempted %v", email, i+1)
	// 		log.Printf("Error: %v", err.Error())

	// 		// exponential backoff
	// 		time.Sleep(time.Second * time.Duration(i+1))
	// 		continue
	// 	}

	// 	log.Printf("Email sent with status code %v", response.StatusCode)
	// 	return nil
	// }

	// return fmt.Errorf("Failed to send email after %v attempts", maxRetries)

	return nil

}
