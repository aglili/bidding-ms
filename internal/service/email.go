package service

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/aglili/auction-app/internal/config"
	"github.com/redis/go-redis/v9"
)

//go:embed templates/*.html
var templatesFS embed.FS

type EmailService struct {
	Config    *config.Config
	QueueKey  string
	Redis     *redis.Client
	templates *template.Template
}

type EmailJob struct {
	To           string `json:"to"`
	Subject      string `json:"subject"`
	TemplateName string `json:"template_name"`
	Data         any    `json:"data"`
}

func NewEmailService(config *config.Config, redis *redis.Client) *EmailService {
	templates, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse email templates: %v", err)
	}

	return &EmailService{
		QueueKey:  "email_queue",
		Redis:     redis,
		Config:    config,
		templates: templates,
	}
}

func (s *EmailService) EnqueueEmail(ctx context.Context, to, subject, templateName string, data any) error {
	job := EmailJob{
		To:           to,
		Subject:      subject,
		TemplateName: templateName,
		Data:         data,
	}

	payload, err := json.Marshal(job)
	if err != nil {
		return err
	}

	if err := s.Redis.LPush(ctx, s.QueueKey, payload).Err(); err != nil {
		return fmt.Errorf("failed to enqueue email job: %w", err)
	}

	return nil
}

func (s *EmailService) sendEmail(job EmailJob) error {
	var body bytes.Buffer

	if err := s.templates.ExecuteTemplate(&body, job.TemplateName, job.Data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", job.TemplateName, err)
	}

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"utf-8\"\r\n\r\n%s",
		s.Config.SmtpSenderEmail, job.To, job.Subject, body.String(),
	))

	fmt.Printf("SMTP Host: %s\n", s.Config.SmtpHost)
	fmt.Printf("SMTP Port: %s\n", s.Config.SmtpPort)
	fmt.Printf("SMTP Email: %s\n", s.Config.SmtpSenderEmail)
	fmt.Printf("SMTP Pass: %s\n", strings.Repeat("*", len(s.Config.SmtpPass)))

	auth := smtp.PlainAuth("", s.Config.SmtpSenderEmail, s.Config.SmtpPass, s.Config.SmtpHost)
	addr := fmt.Sprintf("%s:%s", s.Config.SmtpHost, s.Config.SmtpPort)

	return smtp.SendMail(addr, auth, s.Config.SmtpSenderEmail, []string{job.To}, msg)
}

func (s *EmailService) StartWorker(ctx context.Context) {
	log.Println("Starting email worker...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Email worker stopped")
			return
		default:
			result, err := s.Redis.BRPop(ctx, 5*time.Second, s.QueueKey).Result()
			if err != nil {
				if err.Error() != "redis: nil" {
					log.Printf("Redis error: %v", err)
				}
				continue
			}

			if len(result) < 2 {
				continue
			}

			var job EmailJob
			if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
				log.Printf("Invalid job format: %v", err)
				continue
			}

			if err := s.sendEmail(job); err != nil {
				log.Printf("Failed to send email to %s: %v", job.To, err)
			} else {
				log.Printf("Email sent successfully to %s", job.To)
			}
		}
	}
}