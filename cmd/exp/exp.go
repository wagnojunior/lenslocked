package main

import (
	"fmt"

	"github.com/wagnojunior/lenslocked/models"
)

// These constants are retrieved from the `MailTrap` website
const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 587
	username = "540b454aea2b6e"
	password = "0df11b87e2f0fd"
)

func main() {
	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})

	err := es.ForgotPassword("wagnojunior@gmail.com", "http://some-URL-like-this")
	if err != nil {
		panic(err)
	}
	fmt.Println("email sent")

	// // Header of the email
	// from := "test@lenslocked.com"
	// to := "wagnojunior@gmail.com"
	// subject := "This is a test email"

	// // Body of the email
	// plainText := "This is the body of the email"
	// html := `<h1> Hello there, buddy!</h1><p>This is the email</p><p>Hope you enjoy it!</p>`

	// // Construct the email
	// msg := mail.NewMessage()
	// msg.SetHeader("From", from)
	// msg.SetHeader("To", to)
	// msg.SetHeader("Subject", subject)
	// msg.SetBody("text/plain", plainText)
	// msg.AddAlternative("text/html", html)
	// msg.WriteTo(os.Stdout)

	// // Connect the the mail server and send an email
	// dialer := mail.NewDialer(host, port, username, password)
	// err := dialer.DialAndSend(msg)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Message sent")
}
