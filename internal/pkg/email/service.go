package email

import (
	"fmt"
	"net/smtp"
	"strings"

	"serica-go/internal/conf"
)

type Service struct {
	cfg conf.Email
}

func NewService(cfg *conf.Bootstrap) *Service {
	return &Service{cfg: cfg.Email}
}

func (s *Service) Send(to []string, subject, body string) error {
	if s.cfg.APIKey == "" {
		return fmt.Errorf("email api key not configured")
	}

	from := s.cfg.From
	if from == "" {
		from = "noreply@sericamind.com"
	}

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		from, strings.Join(to, ","), subject, body,
	)

	return smtp.SendMail(
		"smtp.sendgrid.net:587",
		smtp.PlainAuth("apikey", "apikey", s.cfg.APIKey, "smtp.sendgrid.net"),
		from,
		to,
		[]byte(msg),
	)
}

func (s *Service) SendBorrowMail(email, bookTitle, userName, returnDate string) error {
	subject := "金閱閣電子書預約書籍通知"
	body := fmt.Sprintf(`親愛的 %s 讀者:<br/><br/>
你在香港公共圖書館金閱閣電子書預約的《%s》已借入你的帳戶內，期限至 %s。
你可登入香港公共圖書館的 <a href="https://joyread.club">金閱閣電子書</a>，透過瀏覽器線上閱讀。<br/><br/>
多謝使用金閱閣電子書。`, userName, bookTitle, returnDate)
	return s.Send([]string{email}, subject, body)
}

func (s *Service) SendBlockedMail(email, bookTitle, userName string) error {
	subject := "金閱閣電子書預約書籍通知 (借閱冊數已滿)"
	body := fmt.Sprintf(`親愛的 %s 讀者:<br/><br/>
由於你在香港公共圖書館金閱閣電子書的借閱額已滿，你預約的《%s》將排到下一個順位。
請登入香港公共圖書館 <a href="https://joyread.club">金閱閣電子書</a>，歸還任何一本書籍。
待讀者歸還該書後，便會自動借入你的帳戶。<br/><br/>
多謝使用金閱閣電子書。`, userName, bookTitle)
	return s.Send([]string{email}, subject, body)
}
