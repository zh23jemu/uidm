package output

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var (
	LoggerFile    = logrus.New()
	LoggerConsole = logrus.New()
	LogFile       *os.File

	logDir      string
	logMailBody string

	logBodyHead = `<body text='#000000'><center><font size=3 color='#dd0000'><b> UID error logs</b></font></center>
    <br/><table style=' font-size: 12px;'border='1' cellspacing='0' cellpadding='0' bordercolor='#000000' width='95%' align='center' >
    <tr  bgcolor='#B0E0E6' align='left' >
    <th>type</th>
    <th style=width:60px>error</th>
    <th style='width:50px'>details</th>`

	logBodyItem = `<tr align='left' >
    <th style=width:200px>%s</th>
    <th style=width:300px>%s</th>
    <th style=width:500px>%s</th>
    </tr>`
)

func PrepareLog() {
	path, _ := os.Getwd()

	if strings.Contains(path, "\\") { //check OS version
		logDir = path + "\\log\\"
	} else if strings.Contains(path, "/") {
		logDir = path + "/log/"
	}
	logFilename := logDir + "unid-" + time.Now().Format("20060102") + ".log"

	err := os.MkdirAll(logDir, os.ModePerm) //make dir
	if err != nil {
		LoggerConsole.Error(err)
	}

	LogFile, err = os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		LoggerConsole.Error(err)
	}

	LoggerFile.SetFormatter(&log.JSONFormatter{})
	LoggerFile.Out = LogFile

	logMailBody = ""
}

func GenerateLog(err error, typ, details string, isInfo bool) {
	//if subtitle != "" {
	//	title += " - " + subtitle
	//}
	if err != nil {
		LoggerFile.WithFields(logrus.Fields{
			"type":    typ,
			"details": details + "," + err.Error(),
		}).Error()

		logMailBody += fmt.Sprintf(logBodyItem, typ, err.Error(), details)
		//mailbody := fmt.Sprintf("type: %s\ndetails: %s\nerror: %s", typ, details, err.Error())
		//
		//SendMail(MailReport, title, logBody, "", MailAdmin, MailITUserExp)
		//SendMail(MailReport, title, logBodyHead+logMailBody, "html", "", MailAdmin)
	} else if isInfo {
		LoggerConsole.Info(details)
		LoggerFile.WithFields(logrus.Fields{
			"type":    typ,
			"details": details,
		}).Info()
	}
}

func SendLogMail() {
	if logMailBody != "" {
		title := "UNID error log " + time.Now().Format("20060102")
		SendMail(MailReport, title, logBodyHead+logMailBody, "html", "", "", MailAdmin, MailITUserExp)
	}
}
