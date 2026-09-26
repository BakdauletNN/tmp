package notifier

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
)

type SMTPMailer struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSMTPMailer(host, port, username, password, from string) *SMTPMailer {
	return &SMTPMailer{
		host:     host,
		port:     port,
		username: username,
		password: strings.Join(strings.Fields(password), ""),
		from:     from,
	}
}

func (m *SMTPMailer) SendPasswordReset(to, name, resetURL string) error {
	if m.host == "" || m.port == "" || m.from == "" {
		return errors.New("SMTP_HOST, SMTP_PORT, and SMTP_FROM must be configured")
	}
	fromAddress, err := mail.ParseAddress(m.from)
	if err != nil {
		return fmt.Errorf("parse SMTP sender address: %w", err)
	}
	toAddress, err := mail.ParseAddress(to)
	if err != nil {
		return fmt.Errorf("parse recipient address: %w", err)
	}

	auth := smtp.Auth(nil)
	if m.username != "" {
		auth = smtp.PlainAuth("", m.username, m.password, m.host)
	}
	safeName := strings.NewReplacer("\r", " ", "\n", " ").Replace(name)
	subject := "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte("CoworkGo Password Reset")) + "?="
	body := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\nHello, %s.\r\n\r\nTo reset your password, open this link (it is valid for 30 minutes):\r\n%s\r\n\r\nIf you did not request a password reset, simply ignore this email.\r\n",
		fromAddress, toAddress, subject, safeName, resetURL,
	)
	return smtp.SendMail(
		net.JoinHostPort(m.host, m.port),
		auth,
		fromAddress.Address,
		[]string{toAddress.Address},
		[]byte(body),
	)
}
