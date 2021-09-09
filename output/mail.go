package output

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"time"
)

const (
	MailReport    = "service.uid@csisolar.com"
	MailAdmin     = "billy.zhou@csisolar.com"
	MailITUserExp = "F_UserExperience_ALL@csisolar.com"
	MailHelp      = "help@csisolar.com"
	MailOAOp      = "CSI_OA_Operation@csisolar.com"
	MailJoyceZhou = "joyce.zhou@csisolar.com"
	smtpServer    = "10.0.8.27"
	smtpPort      = 25
)

func SendMail(from, subject, body, contentType, attach, bcc string, to ...string) {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	if bcc != "" {
		m.SetHeader("Bcc", bcc)
	}
	//m.SetHeader("Cc", cc)
	m.SetHeader("Subject", subject)
	m.SetBody("text/"+contentType, body)
	if attach != "" {
		m.Attach(attach)
	}

	d := gomail.NewDialer(smtpServer, smtpPort, "", "")

	if err := d.DialAndSend(m); err != nil {
		LoggerFile.WithFields(logrus.Fields{
			"function": "SendMail",
			"type":     "Send mail",
			"detail":   err,
		}).Error(err)
	}
}

func SendMonthlyRemind() {
	startDate := time.Now().Format("0102")
	endDate := time.Now().AddDate(0, 1, 0).Format("0102")

	subject := "月上岗HC预估" + startDate
	body := "Hi Joyce ,\n\n请预估近一个月 (" + startDate + "-" + endDate + ") 集团及各子公司的办公室员工入职的数据。发送到邮箱群组 F_UserExperience_ALL@csisolar.com\n\n谢谢。\n"
	attach := "预估近期入职员工人数（按每月）.xlsx"

	SendMail(MailReport, subject, body, "plain", attach, MailAdmin, MailJoyceZhou)
}
