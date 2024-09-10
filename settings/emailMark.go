package settings

import (
	"Chat/global"
	"Chat/pkg/emailMark"
	"github.com/XYYSWK/Lutils/pkg/email"
)

type mark struct {
}

func (mark) Init() {
	global.EmailMark = emailMark.New(emailMark.Config{
		UserMarkDuration: global.PublicSetting.Rules.UserMarkDuration,
		CodeMarkDuration: global.PublicSetting.Rules.CodeMarkDuration,
		SMTPInfo: email.SMTPInfo{
			Port:     global.PrivateSetting.Email.Port,
			IsSSL:    global.PrivateSetting.Email.IsSSL,
			Host:     global.PrivateSetting.Email.Host,
			UserName: global.PrivateSetting.Email.Username,
			Password: global.PrivateSetting.Email.Password,
			From:     global.PrivateSetting.Email.From,
		},
		AppName: global.PublicSetting.App.Name,
	})
}
