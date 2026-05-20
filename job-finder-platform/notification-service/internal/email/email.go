package email

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type Service struct{ cfg *Config }

func NewService(host, port, username, password, from string) *Service {
	return &Service{cfg: &Config{Host: host, Port: port, Username: username, Password: password, From: from}}
}

func (s *Service) Send(to, subject, body string, isHTML bool) error {
	if s.cfg.Username == "" {
		log.Printf("[EMAIL MOCK] To=%s Subject=%s", to, subject)
		return nil
	}

	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	contentType := "text/plain; charset=UTF-8"
	if isHTML {
		contentType = "text/html; charset=UTF-8"
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "From: JobFinder <%s>\r\n", s.cfg.From)
	fmt.Fprintf(&buf, "To: %s\r\n", to)
	fmt.Fprintf(&buf, "Subject: %s\r\n", subject)
	fmt.Fprintf(&buf, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&buf, "Content-Type: %s\r\n\r\n", contentType)
	buf.WriteString(body)

	addr := s.cfg.Host + ":" + s.cfg.Port
	return smtp.SendMail(addr, auth, s.cfg.From, []string{to}, buf.Bytes())
}

func (s *Service) SendWelcome(to, name, verifyLink string) error {
	body := fmt.Sprintf(`<!DOCTYPE html><html><body style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto">
<div style="background:linear-gradient(135deg,#667eea,#764ba2);padding:30px;border-radius:10px;text-align:center">
  <h1 style="color:white;margin:0">🎉 Welcome to JobFinder!</h1>
</div>
<div style="padding:30px;background:#f9f9f9;border-radius:10px;margin-top:20px">
  <h2>Hello, %s!</h2>
  <p>Thank you for registering. Please verify your email address:</p>
  <div style="text-align:center;margin:30px 0">
    <a href="%s" style="background:#667eea;color:white;padding:15px 30px;border-radius:5px;text-decoration:none;font-size:16px">
      ✅ Verify Email Address
    </a>
  </div>
  <p style="color:#666;font-size:14px">If you didn't create this account, please ignore this email.</p>
</div>
</body></html>`, name, verifyLink)
	return s.Send(to, "Welcome to JobFinder – Verify Your Email", body, true)
}

func (s *Service) SendJobApplied(employerEmail, applicantName, jobTitle, appID string) error {
	body := fmt.Sprintf(`<!DOCTYPE html><html><body style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto">
<div style="background:#2196F3;padding:20px;border-radius:10px">
  <h2 style="color:white;margin:0">📋 New Application Received</h2>
</div>
<div style="padding:20px;background:#f9f9f9;border-radius:10px;margin-top:15px">
  <p>You received a new application for <strong>%s</strong>.</p>
  <p><strong>Applicant:</strong> %s</p>
  <p><strong>Application ID:</strong> %s</p>
  <p>Log in to review this application.</p>
</div>
</body></html>`, jobTitle, applicantName, appID)
	return s.Send(employerEmail, "New Application: "+jobTitle, body, true)
}

func (s *Service) SendStatusUpdate(applicantEmail, applicantName, jobTitle, status, message string) error {
	colors := map[string]string{
		"accepted": "#4CAF50", "rejected": "#f44336",
		"reviewed": "#FF9800", "pending": "#9E9E9E",
	}
	emojis := map[string]string{
		"accepted": "🎉", "rejected": "😔", "reviewed": "👀", "pending": "⏳",
	}
	color := colors[status]
	if color == "" { color = "#9E9E9E" }
	emoji := emojis[status]
	if emoji == "" { emoji = "📋" }

	var msgSection string
	if message != "" {
		msgSection = fmt.Sprintf(`<p><strong>Message:</strong> %s</p>`, message)
	}

	body := fmt.Sprintf(`<!DOCTYPE html><html><body style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto">
<div style="background:%s;padding:20px;border-radius:10px">
  <h2 style="color:white;margin:0">%s Application Update</h2>
</div>
<div style="padding:20px;background:#f9f9f9;border-radius:10px;margin-top:15px">
  <p>Dear <strong>%s</strong>,</p>
  <p>Your application for <strong>%s</strong> has been updated.</p>
  <p><strong>Status:</strong> <span style="color:%s;font-weight:bold;text-transform:uppercase">%s</span></p>
  %s
</div>
</body></html>`, color, emoji, applicantName, jobTitle, color, status, msgSection)
	return s.Send(applicantEmail, fmt.Sprintf("Application Update: %s – %s", jobTitle, strings.Title(status)), body, true)
}

func (s *Service) SendPasswordReset(to, name, resetLink string) error {
	body := fmt.Sprintf(`<!DOCTYPE html><html><body style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto">
<div style="background:#FF5722;padding:20px;border-radius:10px">
  <h2 style="color:white;margin:0">🔒 Password Reset</h2>
</div>
<div style="padding:20px;background:#f9f9f9;border-radius:10px;margin-top:15px">
  <p>Hello <strong>%s</strong>,</p>
  <p>Click below to reset your password (expires in 1 hour):</p>
  <div style="text-align:center;margin:25px 0">
    <a href="%s" style="background:#FF5722;color:white;padding:12px 25px;border-radius:5px;text-decoration:none">
      Reset Password
    </a>
  </div>
  <p style="color:#666;font-size:14px">If you didn't request this, ignore this email.</p>
</div>
</body></html>`, name, resetLink)
	return s.Send(to, "Password Reset – JobFinder", body, true)
}

func (s *Service) SendJobAlert(to, name string, jobTitles []string, query string) error {
	var items strings.Builder
	for _, t := range jobTitles {
		items.WriteString(fmt.Sprintf("<li>%s</li>", t))
	}
	body := fmt.Sprintf(`<!DOCTYPE html><html><body style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto">
<div style="background:#4CAF50;padding:20px;border-radius:10px">
  <h2 style="color:white;margin:0">🔔 New Job Matches</h2>
</div>
<div style="padding:20px;background:#f9f9f9;border-radius:10px;margin-top:15px">
  <p>Hello <strong>%s</strong>,</p>
  <p>New jobs matching <strong>"%s"</strong>:</p>
  <ul>%s</ul>
</div>
</body></html>`, name, query, items.String())
	return s.Send(to, fmt.Sprintf("Job Alert: New jobs matching '%s'", query), body, true)
}
