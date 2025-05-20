package mailer

import (
	"fmt"
	"io/ioutil"
	"net/smtp"
	"path/filepath"

	"github.com/flosch/pongo2"
	"github.com/jordan-wright/email"
	"github.com/vmdt/notification-worker/pkg/logger"
)

type MailerConfig struct {
	SMTPHost     string `mapstructure:"smtp_host"`
	SMTPPort     int    `mapstructure:"smtp_port"`
	SMTPUser     string `mapstructure:"smtp_user"`
	SMTPPassword string `mapstructure:"smtp_password"`
	Sender       string `mapstructure:"sender"`
}

type Mailer struct {
	cfg *MailerConfig
	log logger.ILogger
}

func NewMailer(cfg *MailerConfig, log logger.ILogger) *Mailer {
	return &Mailer{
		cfg: cfg,
		log: log,
	}
}

func (m *Mailer) SendMail(template string, to string, locals map[string]interface{}) error {
	subjectPath := filepath.Join("templates", template, "subject.html")
	htmlPath := filepath.Join("templates", template, "content.html")

	subjectBytes, err := ioutil.ReadFile(subjectPath)
	if err != nil {
		m.log.Errorf("read subject error: %v", err)
		return err
	}

	subjectTpl, err := pongo2.FromString(string(subjectBytes))
	if err != nil {
		m.log.Errorf("parse subject error: %v", err)
		return err
	}

	subject, err := subjectTpl.Execute(locals)
	if err != nil {
		m.log.Errorf("execute subject error: %v", err)
		return err
	}

	// Load HTML template
	htmlBytes, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		m.log.Errorf("read html error: %v", err)
		return err
	}
	htmlTpl, err := pongo2.FromString(string(htmlBytes))
	if err != nil {
		m.log.Errorf("parse html error: %v", err)
		return err
	}
	htmlBody, err := htmlTpl.Execute(locals)
	if err != nil {
		m.log.Errorf("execute html error: %v", err)
		return err
	}

	// Compose email
	e := email.NewEmail()
	e.From = m.cfg.Sender
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(htmlBody)

	addr := fmt.Sprintf("%s:%d", m.cfg.SMTPHost, m.cfg.SMTPPort)
	auth := smtp.PlainAuth("", m.cfg.SMTPUser, m.cfg.SMTPPassword, m.cfg.SMTPHost)

	return e.Send(addr, auth)
}
