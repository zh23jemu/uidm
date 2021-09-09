package common

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	"uidm/output"
)

const (
	timeLayout = "2006-01-02 15:04:05"
)

func ConvertText(text string) string {
	return strings.ReplaceAll(text, "'", "''")
}

func ConvertSFTimeStamp(sfTimeStamp string) string {
	if strings.TrimSpace(sfTimeStamp) == "" {
		return ""
	}
	commonTimeStamp := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(sfTimeStamp, "/Date(", ""), ")/", ""), "+0000", "")

	i, err := strconv.ParseInt(commonTimeStamp, 10, 64)
	output.GenerateLog(err, "Datetime conversion", sfTimeStamp, false)

	if i < -6847833600000 { // The min value allowed of sql server datetime: 1753-01-01 00:00:00
		i = -6847833600000
	}
	return time.Unix(i/1e3, 0).Format(timeLayout)
}

func ConvertDatetimeToDate(dt string) string {
	return strings.Split(dt, "T")[0]
}

func ConvertTitleEN(titleEN string) (title, dept string) {

	if strings.Contains(titleEN, ",") {
		lastComma := strings.LastIndex(titleEN, ",")
		title = strings.TrimSpace(titleEN[:lastComma])
		dept = strings.TrimSpace(titleEN[lastComma+1:])
	} else {
		title = titleEN
	}
	return
}

func ConvertMobile(mobile string) string {
	if len(mobile) == 11 && mobile[0] == '1' {
		var newMobile string
		vv := []rune(mobile)
		for i := 0; i < len(vv); i++ {
			if i == 3 || i == 7 {
				newMobile += " " + string(vv[i])
			} else {
				newMobile += string(vv[i])
			}
		}
		return "+86 " + newMobile
	}
	return mobile
}

func SplitFirstname(text string) []string {
	return strings.Split(strings.TrimSpace(text), " ")
}

func ConvertLastname(text string) string {

	result := strings.ToLower(strings.TrimSpace(text))
	result = strings.ToUpper(result[:1]) + result[1:]

	//var result string
	//vv := []rune(strings.ToUpper(text))
	//for i := 0; i < len(vv); i++ {
	//	if i == 0 {
	//		//vv[i] -= 32 // string的码表相差32位
	//		result += string(vv[i])
	//	} else {
	//		vv[i] += 32
	//		result += string(vv[i])
	//	}
	//}
	return result
}

func ConvertFullname(username string) string {
	if strings.Contains(username, ".") {
		fnames := strings.Split(username, ".")
		fn1 := strings.ToUpper(fnames[0][:1]) + strings.ToLower(fnames[0][1:])
		fn2 := strings.ToUpper(fnames[1][:1]) + strings.ToLower(fnames[1][1:])

		return fn1 + " " + fn2
	} else {
		return username
	}
}

func ConvertDisplayName(text string) string {
	r, _ := regexp.Compile(`[a-zA-Z]+\s[A-Z]+`)
	findResult := r.FindAllString(text, -1)
	if len(findResult) > 0 && findResult[0] == text { // format like "Billy ZHOU"
		fullname := strings.Split(findResult[0], " ")
		//fmt.Println(fullname)
		firstname := fullname[0]
		lastname := fullname[1]

		var upperStr string
		vv := []rune(strings.ToLower(lastname))
		for i := 0; i < len(vv); i++ {
			if i == 0 {
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				upperStr += string(vv[i])
			}
		}
		return firstname + " " + upperStr
	}
	return text
}

func GetStringSliceFirst(s []string) string {
	if len(s) > 0 {
		return s[0]
	} else {
		return ""
	}
}

func ConvertDNToCN(dn string) string {
	if strings.HasPrefix(dn, "CN=") {
		nameWithCN := strings.Split(dn, ",")[0]
		return strings.Split(nameWithCN, "=")[1]
	}
	return ""
}
