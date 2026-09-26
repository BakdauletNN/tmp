package notifier

import (
	"testing"
	"tmp/internal/config"
)

func TestNewSMTPMailerRemovesAppPasswordSeparators(t *testing.T) {
	cfg := config.Config{
		SMTPHost:     "smtp.gmail.com",
		SMTPPort:     "587",
		SMTPUser:     "sender@example.com",
		SMTPPassword: "abcd efgh ijkl mnop",
		SMTPFrom:     "sender@example.com",
	}
	mailer := NewSMTPMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFrom)
	if mailer.password != "abcdefghijklmnop" {
		t.Fatalf("expected grouped app password to be normalized")
	}
}
