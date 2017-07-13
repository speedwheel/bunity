package mail

import (
	"github.com/kataras/go-mailer"
)

type Mail struct {
	mailService mailer.Service
	mailConfig mailer.Config
}


func (this *Mail) Init() *mailer.Service {
	this.mailConfig = mailer.Config{
		Host:     "smtp.sendgrid.net",
		Username: "Speedwheel",
		Password: "logitech11",
		Port:     587,
		FromAddr: "edi.ultras@gmail.com",
		FromAlias: "edi.ultras@gmail.com",
	}

	this.mailService = mailer.New(this.mailConfig)
	return &this.mailService
}

func (this *Mail) Send(to []string, content string) {
	err := this.mailService.Send("iris e-mail just t3st subject", content, to...)

	if err != nil {
		panic(err.Error())
	}
}