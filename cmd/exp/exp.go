package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/wagnojunior/lenslocked/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	host := os.Getenv("SMPT_HOST")
	portStr := os.Getenv("SMPT_PORT")
	port, err := strconv.Atoi(portStr)
	username := os.Getenv("SMPT_USERNAME")
	password := os.Getenv("SMPT_PASSWORD")
	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})

	err = es.ForgotPassword("wagnojunior@gmail.com", "http://some-URL-like-this")
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
