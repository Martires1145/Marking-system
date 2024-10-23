package util

import (
	"crypto/tls"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

var (
	mailer = gomail.NewDialer("smtp.qq.com", 465, viper.GetString("email.account"), viper.GetString("email.password"))
)

func SendMessage(captcha string, email string) error {
	c := "<!DOCTYPE html>\n<html lang=\"en\">\n  <head>\n    <meta charset=\"UTF-8\" />\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\n    <title>验证码</title>\n    <style>\n      body {\n        font-family: Arial, sans-serif;\n        background-color: #f4f4f4;\n        margin: 0;\n        padding: 0;\n        display: flex;\n        justify-content: center;\n        align-items: center;\n        height: 100vh;\n      }\n      .container {\n        background-color: #fff;\n        border-radius: 5px;\n        padding: 20px;\n        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);\n        text-align: center;\n      }\n      h1 {\n        color: #333;\n        margin-bottom: 20px;\n      }\n      p {\n        color: #666;\n        margin-bottom: 20px;\n      }\n      .code {\n        font-size: 24px;\n        color: #007bff;\n        margin-bottom: 20px;\n      }\n      .note {\n        color: #999;\n      }\n    </style>\n  </head>\n  <body>\n    <div class=\"container\">\n      <h1>验证码</h1>\n      <p>您的验证码为：</p>\n      <p class=\"code\">%s</p>\n      <p class=\"note\">请注意，此验证码仅在5分钟内有效。</p>\n    </div>\n  </body>\n</html>"

	content := fmt.Sprintf(c, captcha)
	message := gomail.NewMessage()
	message.SetHeader("From", viper.GetString("email.account")) // 发件人
	message.SetHeader("To", email)                              // 收件人
	message.SetBody("text/html", fmt.Sprintf(content))
	message.SetHeader("Subject", "验证码")

	mailer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := mailer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
