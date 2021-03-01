package tools

import (
	"bytes"
	"net/smtp"

	"github.com/jordan-wright/email"
)

var (
	//EmailFromUser 以下的全局变量暂定
	EmailFromUser = "1015138659@qq.com"
	//EmailToUser 以下的全局变量暂定
	EmailToUser = "m17628294021@163.com"
	//EmailContent 以下的全局变量暂定
	EmailContent = "新版框架-你有新消息注意查收"
	//EmailPassword 以下的全局变量暂定
	EmailPassword = "ntlufpojrqetbehe"
	//EmailHost 以下的全局变量暂定
	EmailHost = "smtp.qq.com"
	//EmailAddr 以下的全局变量暂定
	EmailAddr = "smtp.qq.com:587"
)

//SendMail 发送邮件
func SendMail(fromUser, toUser, subject, attachFile string) error {
	// NewEmail返回一个email结构体的指针
	e := email.NewEmail()
	// 发件人
	e.From = fromUser
	// 收件人(可以有多个)
	e.To = []string{toUser}
	// 邮件主题
	e.Subject = subject

	body := new(bytes.Buffer)
	body.Write([]byte("这里填邮件内容"))

	e.HTML = body.Bytes()
	// 从缓冲中将内容作为附件到邮件中
	if attachFile != "" {
		e.AttachFile(attachFile)
	}

	// 以路径将文件作为附件添加到邮件中
	e.AttachFile("/home/shuai/go/src/email/main.go")
	//fmt.Println("EmailAddr : ",ConfigJson.EmailAddr," EmailFromUser: ",ConfigJson.EmailFromUser)
	//fmt.Println("EmailPassword : ",ConfigJson.EmailPassword," EmailHost: ",ConfigJson.EmailHost)
	// 发送邮件(如果使用QQ邮箱发送邮件的话，passwd不是邮箱密码而是授权码)
	return e.Send(EmailAddr, smtp.PlainAuth("", fromUser, EmailPassword, EmailHost))
	//return e.Send("smtp.qq.com:587", smtp.PlainAuth("", "1015138659@qq.com", "quaybszytcjcbgaj", "smtp.qq.com"))
}
