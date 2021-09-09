package output

import (
	"fmt"
	"net/http"
	"strings"
)

func CreateTicket(url, key, data string) {
	var r http.Request

	r.ParseForm()
	r.Form.Add("input_data", data)
	bodystr := strings.TrimSpace(r.Form.Encode())
	req, err := http.NewRequest("POST", url, strings.NewReader(bodystr))
	if err != nil {
		fmt.Println(err.Error())
		GenerateLog(err, "Wrapping HTTP request", url+" | "+bodystr, false)
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("TECHNICIAN_KEY", key)
		req.Header.Set("Connection", "Keep-Alive")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	GenerateLog(err, "Sending HTTP request", url+" | "+bodystr, false)
	if err == nil {
		resp.Body.Close()
	}else{
		fmt.Println(err.Error())
	}
}
