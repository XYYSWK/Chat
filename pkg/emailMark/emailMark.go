package emailMark

import (
	"errors"
	"fmt"
	"github.com/XYYSWK/Rutils/pkg/email"
	"sync"
	"time"
)

/*
邮箱验证码标记
*/

type EmailMark struct {
	config   Config
	userMark sync.Map //标记用户
	codeMark sync.Map //记录 code
}

type Config struct {
	UserMarkDuration time.Duration  // 用户标记时长
	CodeMarkDuration time.Duration  // 验证码标记时长
	SMTPInfo         email.SMTPInfo //邮箱配置
	AppName          string         //应用名称
}

func New(config Config) *EmailMark {
	return &EmailMark{
		config:   config,
		userMark: sync.Map{},
		codeMark: sync.Map{},
	}
}

var ErrSendTooMany = errors.New("发送过于频繁")

// CheckUserExist 判断邮箱是否已经被记录
func (m *EmailMark) CheckUserExist(emailStr string) bool {
	//判断该键是否存在
	_, ok := m.userMark.Load(emailStr)
	return ok
}

// SendMark 发送验证码
func (m *EmailMark) SendMark(emailStr, code string) error {
	// 发送频率限制
	if m.CheckUserExist(emailStr) {
		return ErrSendTooMany
	}
	m.userMark.Store(emailStr, struct{}{})
	sendMark := email.NewEmail(&m.config.SMTPInfo)
	// 发送邮件
	err := sendMark.SendMail([]string{emailStr}, fmt.Sprintf("%s邮箱验证码", m.config.AppName), fmt.Sprintf("<h1>邮箱验证码</h1>尊敬的用户您好！<br>您的验证码是：%s，请在 %v 分钟内进行验证。O(∩_∩)~", code, m.config.CodeMarkDuration.Minutes()))
	if err != nil {
		// 发送失败删除标记
		m.userMark.Delete(emailStr)
		return err
	}
	// 记录验证码
	m.codeMark.Store(emailStr, code)
	// 验证码过期或用户一定间隔时间后被允许再次请求验证码，延时删除标记
	m.DeleteMarkDelay(emailStr)
	return nil
}

// DeleteMarkDelay 验证码过期，延时删除标记
func (m *EmailMark) DeleteMarkDelay(emailStr string) {
	time.AfterFunc(m.config.CodeMarkDuration, func() {
		m.codeMark.Delete(emailStr)
	})
	time.AfterFunc(m.config.UserMarkDuration, func() {
		m.userMark.Delete(emailStr)
	})
}

// CheckCode 校验验证码
func (m *EmailMark) CheckCode(emailStr, code string) bool {
	myCode, ok := m.codeMark.Load(emailStr)
	// 验证成功删除标记
	if ok && myCode == code {
		m.codeMark.Delete(emailStr)
		return true
	}
	return false
}
