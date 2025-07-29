package mail

import (
	"bytes"
	"fmt"
	"homestay-be/cmd/types"
	"html/template"

	"gopkg.in/gomail.v2"
)

type MailClient struct {
	from     string
	username string
	password string
	host     string
	port     int
}

// NewMailClient khởi tạo client với thông tin SMTP
func NewMailClient(from, username, password string) *MailClient {
	return &MailClient{
		from:     from,
		username: username,
		password: password,
		host:     "in-v3.mailjet.com", // Mailjet SMTP host
		port:     587,
	}
}

// SendText gửi email dạng text
func (mc *MailClient) SendText(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mc.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(mc.host, mc.port, mc.username, mc.password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("gửi email thất bại: %w", err)
	}

	return nil
}


// SendBookingConfirmation gửi email xác nhận đặt phòng với HTML template
func (mc *MailClient) SendBookingConfirmation(to string, data types.BookingEmailData) error {
	tmpl, err := template.ParseFiles("cmd/templates/confirm_email.html")
	if err != nil {
		return fmt.Errorf("không thể load template email: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("render HTML thất bại: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", mc.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", fmt.Sprintf("Xác nhận đặt phòng tại %s", data.HomestayName))
	m.SetBody("text/html", buf.String())

	d := gomail.NewDialer(mc.host, mc.port, mc.username, mc.password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("gửi email thất bại: %w", err)
	}

	return nil
}