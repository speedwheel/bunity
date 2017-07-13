package fb

import (
	"github.com/huandu/facebook"
)

type Fb struct {
	Api *facebook.App
	Session *facebook.Session
	Params facebook.Params
}

func (this *Fb) Init(token string) (*facebook.Session, error) {
	this.Api = facebook.New("310751436051153", "7570aabfe680cdff53179190869134fb")
	this.Session = this.Api.Session(token)
	err := this.Session.Validate()
	if err != nil {
		panic(err.Error())
	}
	this.Params = facebook.Params{
		"fields": "first_name,last_name,email,id,picture.type(large)",
	}
	return this.Session, err
}

func (this *Fb) GetUser() *facebook.Result {
	res, _ := this.Session.Get("/me", this.Params)
	
	return &res
}
