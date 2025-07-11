package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	redispkg "github.com/go-redis/redis/v8"
	"github.com/msyaifudin/pos/internal/redis"
	"gopkg.in/gomail.v2"
)

// EmailJob represents a job for the email worker
type EmailJob struct {
	To      string
	Subject string
	Body    string
}

// EmailQueue is a channel that acts as a queue for email jobs
var EmailQueue chan EmailJob

// InitEmailQueue initializes the email queue
func InitEmailQueue() {
	// A buffered channel to handle up to 100 emails at a time.
	EmailQueue = make(chan EmailJob, 100)
}

// StartEmailWorker starts a worker that processes email jobs from the queue
func StartEmailWorker() {
	go func() {
		for job := range EmailQueue {
			if err := sendEmail(job.To, job.Subject, job.Body); err != nil { // Revert call
				log.Printf("Failed to send email to %s: %v", job.To, err)
			}
		}
	}()
}

// sendEmail contains the logic to send an email
func sendEmail(to, subject, body string) error { // Revert signature
	mailHost := os.Getenv("MAIL_HOST")
	mailPortStr := os.Getenv("MAIL_PORT")
	mailUsername := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")
	mailFromAddress := os.Getenv("MAIL_FROM_ADDRESS")
	mailFromName := os.Getenv("MAIL_FROM_NAME")

	mailPort, err := strconv.Atoi(mailPortStr)
	if err != nil {
		log.Printf("Invalid MAIL_PORT: %v", err)
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", mailFromName, mailFromAddress))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(mailHost, mailPort, mailUsername, mailPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Could not send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

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
	// Check if email can be sent based on rate limit
	if !CanSendEmail(to) {
		log.Printf("Email to %s rate limited. Please wait before sending another email.", to)
		return fmt.Errorf("email rate limited")
	}

	// Read the HTML template from the file
	templateBytes, err := os.ReadFile("internal/templates/emails/email_template.html")
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
		OTP     string
		LogoURL string
	}{
		OTP: otp,
	}

	// Get HOST from environment variable
	logo := os.Getenv("LOGO")
	if logo == "" {
		log.Println("Warning: HOST environment variable not set. Logo URL might be incomplete.")
	}

	// Construct logo URL
	data.LogoURL = logo

	// Execute the template with the data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Printf("Could not execute email template: %v", err)
		return err
	}

	// Create a job and add it to the queue
	job := EmailJob{
		To:      to,
		Subject: "Verification Code",
		Body:    body.String(),
	}

	EmailQueue <- job
	log.Printf("Email job for %s queued.", to)

	return nil
}

// CanSendEmail checks if an email can be sent to a recipient based on a 1-minute cooldown
func CanSendEmail(email string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("email_cooldown:%s", email)

	// Check if the key exists in Redis
	val, err := redis.Rdb.Get(ctx, key).Result()
	if err == redispkg.Nil {
		// Key does not exist, so we can send the email. Set the key with a 1-minute expiry.
		redis.Rdb.Set(ctx, key, "1", 1*time.Minute)
		return true
	} else if err != nil {
		// An error occurred with Redis, log it and allow sending to avoid blocking legitimate emails
		log.Printf("Redis error checking email cooldown for %s: %v", email, err)
		return true
	}

	// Key exists, meaning an email was sent recently. Do not send.
	log.Printf("Email to %s is on cooldown. Last sent at %s", email, val)
	return false
}
