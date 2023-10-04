package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"

	"github.com/brianbianchi/rss/util"
	"github.com/mmcdole/gofeed"
)

var (
	fromEmail    = os.Getenv("FROM")
	smtpUsername = os.Getenv("SMTP_USER")
	smtpPassword = os.Getenv("SMTP_PASS")
	smtpPort     = os.Getenv("SMTP_PORT")
	smtpServer   = os.Getenv("SMTP_SERVER")
)

type EmailData struct {
	BaseUrl string
	Code    string
	Rss     []*gofeed.Feed
}

func SendSubEmail(to string, rss []*gofeed.Feed, code string) {
	body := prepareRssForEmail(rss, code)
	subject := "Daily subscritions"

	message := fmt.Sprintf("From: %s\r\n", fromEmail)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-version: 1.0;\r\n"
	message += "Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	message += fmt.Sprintf("%s\r\n", body)

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, fromEmail, []string{to}, []byte(message))
	if err != nil {
		fmt.Println(err)
	}
}

func prepareRssForEmail(rss []*gofeed.Feed, code string) string {
	path := util.GetRootPath()
	t, err := template.ParseFiles(fmt.Sprint(path, "mail/subs.html"))
	if err != nil {
		panic(err)
	}

	emailData := EmailData{
		BaseUrl: os.Getenv("BASEURL"),
		Code:    code,
		Rss:     rss,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, emailData); err != nil {
		panic(err)
	}

	return tpl.String()
}
