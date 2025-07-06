package services

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)



func GenerateOTP() (string, error) {
	const otpChars = "0123456789"
	const otpLength = 6

	buffer := make([]byte, otpLength)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < otpLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(otpCharsLength)))
		if err != nil {
			return "", err
		}
		buffer[i] = otpChars[num.Int64()]
	}

	return string(buffer), nil
}

func SendVerificationEmail(to, otp string) error {
	// Read the HTML template from the file
	templateBytes, err := os.ReadFile("internal/services/email_template.html")
	if err != nil {
		log.Printf("Could not read email template: %v", err)
		return err
	}

	// Parse the template
	tmpl, err := template.New("emailTemplate").Parse(string(templateBytes))
	if err != nil {
		log.Printf("Could not parse email template: %v", err)
		return err
	}

	// Create a data struct to pass to the template
	data := struct {
		OTP string
	}{
		OTP: otp,
	}

	// Execute the template with the data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Printf("Could not execute email template: %v", err)
		return err
	}

	mailHost := os.Getenv("MAIL_HOST")
	mailPort := os.Getenv("MAIL_PORT")
	mailUsername := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")
	mailFromAddress := os.Getenv("MAIL_FROM_ADDRESS")
	mailFromName := os.Getenv("MAIL_FROM_NAME")

	port, err := strconv.Atoi(mailPort)
	if err != nil {
		log.Printf("Invalid MAIL_PORT: %v", err)
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", mailFromName, mailFromAddress))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Verification Code")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(mailHost, port, mailUsername, mailPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Could not send email: %v", err)
		return err
	}

	return nil
}
