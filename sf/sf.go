package sf

import (
	"encoding/json"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"uidm/common"
	"uidm/db"
	"uidm/output"
)

func GetSFJsonRaw(api string) map[string]json.RawMessage {
	var jsonRaw map[string]json.RawMessage
	client := &http.Client{}

	req, err := http.NewRequest("GET", api, nil)
	output.GenerateLog(err, "SF API http request", api,  false)

	req.Header.Add("Authorization", sfAuth)
	resp, err := client.Do(req)
	output.GenerateLog(err, "SF API authorization", sfAuth, false)

	body, err := ioutil.ReadAll(resp.Body)
	output.GenerateLog(err, "HTTP read", "", false)

	jsonStr := string(body)

	err = json.Unmarshal([]byte(jsonStr), &jsonRaw) // convert to raw json
	output.GenerateLog(err, "Raw json parsing", jsonStr, false)

	return jsonRaw
}

func SyncPerPersonal(hours int) {
	var (
		dataSlice []PerPersonal
	)
	sfApi := mapSFApi["PerPersonal"] //&$filter=personIdExternal eq 'CSVN0469'" //&customPageSize=200"

	if hours > 0 { // update operation
		h, _ := time.ParseDuration("-" + strconv.Itoa(hours) + "h")
		sfApi += "&$filter=lastModifiedDateTime gt datetime'" + time.Now().Add(h).Format("2006-01-02T15:04:05") + "'" // + "T00:00:00'"
	}
	sfApiNext := sfApi

	//sfAPIUrl := "https://api15.sapsf.cn:443/odata/v2/PerPersonal?$format=json&$skip=59410&customPageSize=10"
	//sfAPIUrl = "https://api15.sapsf.cn:443/odata/v2/PerPersonal?$format=json&$skip=57000"
	//sfAPIUrl = "https://api15.sapsf.cn:443/odata/v2/PerPersonal?$format=JSON&$orderby=personIdExternal&$filter=lastModifiedOn gt datetime'2021-01-15T00:00:00'"

	for { // query data from SF, then insert to struct slice
		jsonRaw := GetSFJsonRaw(sfApiNext)

		dataResults := new(PerPersonalResults)            // must be renewed here, or will use the last value when meets the null value
		err := json.Unmarshal(jsonRaw["d"], &dataResults) // remove the first "d", and put result to EmployeeResult struct
		output.GenerateLog(err, "Results json parsing", string(jsonRaw["d"]), false)

		dataSlice = append(dataSlice, dataResults.Results...)

		sfApiNext = dataResults.NextURL
		if sfApiNext == "" { //here we meet the end of the query
			break
		}
	}
	output.GenerateLog(nil, "SF data retrieving", strconv.Itoa(len(dataSlice))+" records retrieved in "+sfApi, true)

	b := pgb.NewBar("PerPersonal", len(dataSlice))
	for _, data := range dataSlice {
		b.Add()

		personIDExternal := data.PersonIDExternal
		startDate := common.ConvertSFTimeStamp(data.StartDate)
		lastName := strings.ReplaceAll(data.LastName, "'", "''")
		lastModifiedDateTime := common.ConvertSFTimeStamp(data.LastModifiedDateTime)
		gender := data.Gender
		endDate := common.ConvertSFTimeStamp(data.EndDate)
		createdOn := common.ConvertSFTimeStamp(data.CreatedOn)
		customString5 := data.CustomString5
		customString6 := strings.ReplaceAll(data.CustomString6, "'", "''")
		firstName := strings.ReplaceAll(data.FirstName, "'", "''")

		sqlquery := fmt.Sprintf(`IF EXISTS(SELECT personIdExternal FROM PerPersonal WHERE personIdExternal = '%s') UPDATE PerPersonal SET [startDate]='%s', [lastName]=N'%s', [lastModifiedDateTime]='%s', [gender]='%s', [endDate]='%s', [createdOn]='%s', [customString5]='%s', [customString6]='%s', [firstName]=N'%s' WHERE [personIdExternal] ='%s' ELSE INSERT INTO PerPersonal (personIdExternal,startDate,lastName,lastModifiedDateTime,gender,endDate,createdOn,customString5,customString6,firstName) VALUES('%s','%s',N'%s','%s','%s','%s','%s','%s','%s',N'%s');`, personIDExternal, startDate, lastName, lastModifiedDateTime, gender, endDate, createdOn, customString5, customString6, firstName, personIDExternal, personIDExternal, startDate, lastName, lastModifiedDateTime, gender, endDate, createdOn, customString5, customString6, firstName)
		db.QueryMsSql(sqlquery)
	}
	fmt.Println()
}

func SyncFOCompany(hours int) {
	var (
		dataSlice []FOCompany
	)
	sfApi := mapSFApi["FOCompany"] //&$filter=personIdExternal eq 'CSVN0469'" //&customPageSize=200"

	if hours > 0 { // update operation
		h, _ := time.ParseDuration("-" + strconv.Itoa(hours) + "h")
		sfApi += "&$filter=lastModifiedDateTime gt datetime'" + time.Now().Add(h).Format("2006-01-02T15:04:05") + "'" // + "T00:00:00'"
	}
	sfApiNext := sfApi

	for { // query data from SF, then insert to struct slice
		jsonRaw := GetSFJsonRaw(sfApiNext)

		dataResults := new(FOCompanyResults)              // must be renewed here, or will use the last value when meets the null value
		err := json.Unmarshal(jsonRaw["d"], &dataResults) // remove the first "d", and put result to EmployeeResult struct
		output.GenerateLog(err, "Results json parsing", string(jsonRaw["d"]), false)

		dataSlice = append(dataSlice, dataResults.Results...)

		sfApiNext = dataResults.NextURL
		if sfApiNext == "" { //here we meet the end of the query
			break
		}
	}
	output.GenerateLog(nil, "SF data retrieving", strconv.Itoa(len(dataSlice))+" records retrieved in "+sfApi, true)

	b := pgb.NewBar("FOCompany", len(dataSlice))
	for _, data := range dataSlice {
		b.Add()

		externalCode := data.ExternalCode
		startDate := common.ConvertSFTimeStamp(data.StartDate)
		country := data.Country
		nameLocalized := data.NameLocalized
		nameThTH := data.NameThTH
		name := data.Name
		status := data.Status
		description := data.Description
		endDate := common.ConvertSFTimeStamp(data.EndDate)
		createdOn := common.ConvertSFTimeStamp(data.CreatedOn)
		lastModifiedDateTime := common.ConvertSFTimeStamp(data.LastModifiedDateTime)
		nameEnUS := data.NameEnUS
		currency := data.Currency
		nameZhCN := data.NameZhCN
		defaultLocation := data.DefaultLocation
		nameDefaultValue := data.NameDefaultValue

		sqlquery := fmt.Sprintf(`IF EXISTS(SELECT externalCode FROM FOCompany WHERE externalCode = '%s') UPDATE FOCompany SET [startDate]='%s', [country]='%s', [name_localized]='%s', [name_th_TH]='%s', [name]='%s', [status]='%s', [description]='%s', [endDate]='%s', [createdOn]='%s', [lastModifiedDateTime]='%s', [name_en_US]='%s', [currency]='%s', [name_zh_CN]='%s', [defaultLocation]='%s', [name_defaultValue]='%s' WHERE [externalCode] ='%s' ELSE INSERT INTO FOCompany (externalCode,startDate,country,name_localized,name_th_TH,name,status,description,endDate,createdOn,lastModifiedDateTime,name_en_US,currency,name_zh_CN,defaultLocation,name_defaultValue) VALUES('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s');`, externalCode, startDate, country, nameLocalized, nameThTH, name, status, description, endDate, createdOn, lastModifiedDateTime, nameEnUS, currency, nameZhCN, defaultLocation, nameDefaultValue, externalCode, externalCode, startDate, country, nameLocalized, nameThTH, name, status, description, endDate, createdOn, lastModifiedDateTime, nameEnUS, currency, nameZhCN, defaultLocation, nameDefaultValue)
		db.QueryMsSql(sqlquery)
	}
	fmt.Println()
}

func SyncEmpEmployment(hours int) {
	var (
		dataSlice []EmpEmployment
	)
	sfApi := mapSFApi["EmpEmployment"] //&$filter=personIdExternal eq 'CSVN0469'" //&customPageSize=200"

	if hours > 0 { // update operation
		h, _ := time.ParseDuration("-" + strconv.Itoa(hours) + "h")
		sfApi += "&$filter=lastModifiedDateTime gt datetime'" + time.Now().Add(h).Format("2006-01-02T15:04:05") + "'" // + "T00:00:00'"
	}
	sfApiNext := sfApi

	for { // query data from SF, then insert to struct slice
		jsonRaw := GetSFJsonRaw(sfApiNext)

		dataResults := new(EmpEmploymentResults)          // must be renewed here, or will use the last value when meets the null value
		err := json.Unmarshal(jsonRaw["d"], &dataResults) // remove the first "d", and put result to EmployeeResult struct
		output.GenerateLog(err, "Results json parsing", string(jsonRaw["d"]), false)

		dataSlice = append(dataSlice, dataResults.Results...)

		sfApiNext = dataResults.NextURL
		if sfApiNext == "" { //here we meet the end of the query
			break
		}
	}
	output.GenerateLog(nil, "SF data retrieving", strconv.Itoa(len(dataSlice))+" records retrieved in "+sfApi, true)

	b := pgb.NewBar("EmpEmployment", len(dataSlice))
	for _, data := range dataSlice {
		b.Add()

		personIDExternal := data.PersonIDExternal
		userID := data.UserID
		endDate := common.ConvertSFTimeStamp(data.EndDate)
		assignmentIDExternal := data.AssignmentIDExternal
		seniorityDate := common.ConvertSFTimeStamp(data.SeniorityDate)
		startDate := common.ConvertSFTimeStamp(data.StartDate)
		hiringNotCompleted := strconv.FormatBool(data.HiringNotCompleted)
		isEcRecord := strconv.FormatBool(data.IsECRecord)
		customString20 := data.CustomString20
		createOn := common.ConvertSFTimeStamp(data.CreatedOn)
		lastDateWorked := common.ConvertSFTimeStamp(data.LastDateWorked)
		originalStartDate := common.ConvertSFTimeStamp(data.OriginalStartDate)
		assignmentClass := data.AssignmentClass
		customString5 := data.CustomString5
		lastModifiedDateTime := common.ConvertSFTimeStamp(data.LastModifiedDateTime)

		sqlquery := fmt.Sprintf(`IF EXISTS(SELECT personIdExternal FROM EmpEmployment WHERE personIdExternal = '%s') UPDATE EmpEmployment SET [userId]='%s', [endDate]='%s', [assignmentIdExternal]='%s', [seniorityDate]='%s', [startDate]='%s', [hiringNotCompleted]='%s', [isECRecord]='%s', [customString20]='%s', [createdOn]='%s', [lastDateWorked]='%s', [originalStartDate]='%s', [assignmentClass]='%s', [customString5]='%s', [lastModifiedDateTime]='%s' WHERE [personIdExternal] ='%s' ELSE INSERT INTO EmpEmployment (personIdExternal,userId,endDate,assignmentIdExternal,seniorityDate,startDate,hiringNotCompleted,isECRecord,customString20,createdOn,lastDateWorked,originalStartDate,assignmentClass,customString5,lastModifiedDateTime) VALUES('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s');`, personIDExternal, userID, endDate, assignmentIDExternal, seniorityDate, startDate, hiringNotCompleted, isEcRecord, customString20, createOn, lastDateWorked, originalStartDate, assignmentClass, customString5, lastModifiedDateTime, personIDExternal, personIDExternal, userID, endDate, assignmentIDExternal, seniorityDate, startDate, hiringNotCompleted, isEcRecord, customString20, createOn, lastDateWorked, originalStartDate, assignmentClass, customString5, lastModifiedDateTime)
		db.QueryMsSql(sqlquery)
	}
	fmt.Println()
}

func SyncEmpJob(hours int) {
	var (
		dataSlice []EmpJob
	)
	sfApi := mapSFApi["EmpJob"] //&$filter=personIdExternal eq 'CSVN0469'" //&customPageSize=200"

	if hours > 0 { // update operation
		h, _ := time.ParseDuration("-" + strconv.Itoa(hours) + "h")
		sfApi += "&$filter=lastModifiedDateTime gt datetime'" + time.Now().Add(h).Format("2006-01-02T15:04:05") + "'" // + "T00:00:00'"
	}
	//sfAPIUrl = "https://api15.sapsf.cn:443/odata/v2/EmpJob?$format=JSON&$orderby=userId&$filter=userId eq 'S2E03566'"
	sfApiNext := sfApi

	for { // query data from SF, then insert to struct slice
		jsonRaw := GetSFJsonRaw(sfApiNext)

		dataResults := new(EmpJobResults)                 // must be renewed here, or will use the last value when meets the null value
		err := json.Unmarshal(jsonRaw["d"], &dataResults) // remove the first "d", and put result to EmployeeResult struct
		output.GenerateLog(err, "Results json parsing", string(jsonRaw["d"]), false)

		dataSlice = append(dataSlice, dataResults.Results...)

		sfApiNext = dataResults.NextURL
		if sfApiNext == "" { //here we meet the end of the query
			break
		}
	}
	output.GenerateLog(nil, "SF data retrieving", strconv.Itoa(len(dataSlice))+" records retrieved in "+sfApi, true)

	b := pgb.NewBar("EmpJob", len(dataSlice))
	for _, data := range dataSlice {
		b.Add()

		userID := data.UserID
		startDate := common.ConvertSFTimeStamp(data.StartDate)
		jobCode := data.JobCode
		division := data.Division
		emplStatus := data.EmplStatus
		countryOfCompany := data.CountryOfCompany
		managerID := data.ManagerID
		lastModifiedDateTime := common.ConvertSFTimeStamp(data.LastModifiedDateTime)
		createOn := common.ConvertSFTimeStamp(data.CreatedOn)
		position := data.Position
		company := data.Company
		department := data.Department
		employeeClass := data.EmployeeClass
		location := data.Location
		jobTitle := data.JobTitle
		customString10 := data.CustomString10

		sqlquery := fmt.Sprintf(`IF EXISTS(SELECT userId FROM EmpJob WHERE userId = '%s') UPDATE EmpJob SET [startDate]='%s', [jobCode]='%s', [division]='%s', [emplStatus]='%s', [countryOfCompany]='%s', [managerId]='%s', [lastModifiedDateTime]='%s', [createdOn]='%s', [position]='%s', [company]='%s', [department]='%s', [employeeClass]='%s', [location]='%s' , [jobTitle]='%s', [customString10]='%s' WHERE [userId] ='%s' ELSE INSERT INTO EmpJob (userId,startDate,jobCode,division,emplStatus,countryOfCompany,managerId,lastModifiedDateTime,createdOn,position,company,department,employeeClass,location,jobTitle,customString10) VALUES('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s');`, userID, startDate, jobCode, division, emplStatus, countryOfCompany, managerID, lastModifiedDateTime, createOn, position, company, department, employeeClass, location, jobTitle, customString10, userID, userID, startDate, jobCode, division, emplStatus, countryOfCompany, managerID, lastModifiedDateTime, createOn, position, company, department, employeeClass, location, jobTitle, customString10)

		db.QueryMsSql(sqlquery)
	}
	fmt.Println()
}

func SyncFODepartment(hours int) {
	var (
		dataSlice []FODepartment
	)
	sfApi := mapSFApi["FODepartment"] //sfApi += `&$filter=externalCode eq '00003812'`

	if hours > 0 { // update operation
		h, _ := time.ParseDuration("-" + strconv.Itoa(hours) + "h")
		sfApi += "&$filter=lastModifiedDateTime gt datetime'" + time.Now().Add(h).Format("2006-01-02T15:04:05") + "'" // + "T00:00:00'"
	}
	sfApiNext := sfApi

	for { // query data from SF, then insert to struct slice
		jsonRaw := GetSFJsonRaw(sfApiNext)

		dataResults := new(FODepartmentResults)           // must be renewed here, or will use the last value when meets the null value
		err := json.Unmarshal(jsonRaw["d"], &dataResults) // remove the first "d", and put result to EmployeeResult struct
		output.GenerateLog(err, "Results json parsing", string(jsonRaw["d"]), false)

		dataSlice = append(dataSlice, dataResults.Results...)

		sfApiNext = dataResults.NextURL
		if sfApiNext == "" { //here we meet the end of the query
			break
		}
	}
	output.GenerateLog(nil, "SF data retrieving", strconv.Itoa(len(dataSlice))+" records retrieved in "+sfApi, true)

	b := pgb.NewBar("FODepartment", len(dataSlice))
	for _, data := range dataSlice {
		b.Add()

		externalCode := common.ConvertText(data.ExternalCode)
		startDate := common.ConvertSFTimeStamp(data.StartDate)
		parent := common.ConvertText(data.Parent)
		nameLocalized := common.ConvertText(data.NameLocalized)
		custDeptLevel := common.ConvertText(data.CustDeptLevel)
		descriptionDefaultValue := common.ConvertText(data.DescriptionDefaultValue)
		name := common.ConvertText(data.Name)
		descriptionEnUS := common.ConvertText(data.DescriptionEnUS)
		status := common.ConvertText(data.Status)
		description := common.ConvertText(data.Description)
		createOn := common.ConvertSFTimeStamp(data.CreatedOn)
		nameEnUS := common.ConvertText(data.NameEnUS)
		lastModifiedDateTime := common.ConvertSFTimeStamp(data.LastModifiedDateTime)
		descriptionLocalized := common.ConvertText(data.DescriptionLocalized)
		custHeadPosition := common.ConvertText(data.CustHeadPosition)
		headOfUnit := common.ConvertText(data.HeadOfUnit)

		sqlquery := fmt.Sprintf(`IF EXISTS(SELECT externalCode FROM FODepartment WHERE externalCode = '%s') UPDATE FODepartment SET [startDate]='%s', [parent]='%s', [name_localized]='%s', [cust_deptLevel]='%s', [description_defaultValue]=N'%s', [name]='%s', [description_en_US]=N'%s', [status]='%s', [description]=N'%s', [createdOn]='%s', [name_en_US]='%s', [lastModifiedDateTime]='%s', [description_localized]=N'%s', [cust_headPosition]='%s', [headOfUnit]='%s' WHERE [externalCode] ='%s' ELSE INSERT INTO FODepartment (externalCode,startDate,parent,name_localized,cust_deptLevel,description_defaultValue,name,description_en_US,status,description,createdOn,name_en_US,lastModifiedDateTime,description_localized,cust_headPosition,headOfUnit) VALUES('%s','%s','%s','%s','%s',N'%s','%s',N'%s','%s',N'%s','%s','%s','%s',N'%s','%s','%s');`, externalCode, startDate, parent, nameLocalized, custDeptLevel, descriptionDefaultValue, name, descriptionEnUS, status, description, createOn, nameEnUS, lastModifiedDateTime, descriptionLocalized, custHeadPosition, headOfUnit, externalCode, externalCode, startDate, parent, nameLocalized, custDeptLevel, descriptionDefaultValue, name, descriptionEnUS, status, description, createOn, nameEnUS, lastModifiedDateTime, descriptionLocalized, custHeadPosition, headOfUnit)

		db.QueryMsSql(sqlquery)
	}
	fmt.Println()
}

func SyncPerEmail(hours int) {
	var (
		dataSlice []PerEmail
	)
	sfApi := mapSFApi["PerEmail"] //sfApi += `&$filter=externalCode eq '00003812'`

	if hours > 0 { // update operation
		h, _ := time.ParseDuration("-" + strconv.Itoa(hours) + "h")
		sfApi += "&$filter=lastModifiedDateTime gt datetime'" + time.Now().Add(h).Format("2006-01-02T15:04:05") + "'" // + "T00:00:00'"
	}
	sfApiNext := sfApi

	for { // query data from SF, then insert to struct slice
		jsonRaw := GetSFJsonRaw(sfApiNext)

		dataResults := new(PerEmailResults)               // must be renewed here, or will use the last value when meets the null value
		err := json.Unmarshal(jsonRaw["d"], &dataResults) // remove the first "d", and put result to EmployeeResult struct
		output.GenerateLog(err, "Results json parsing", string(jsonRaw["d"]), false)

		dataSlice = append(dataSlice, dataResults.Results...)

		sfApiNext = dataResults.NextURL
		if sfApiNext == "" { //here we meet the end of the query
			break
		}
	}
	output.GenerateLog(nil, "SF data retrieving", strconv.Itoa(len(dataSlice))+" records retrieved in "+sfApi, true)

	b := pgb.NewBar("PerMail", len(dataSlice))
	for _, data := range dataSlice {
		b.Add()

		emailType := data.EmailType
		personIDExternal := data.PersonIDExternal
		emailAddress := data.EmailAddress
		createdOn := common.ConvertSFTimeStamp(data.CreatedOn)
		lastModifiedDateTime := common.ConvertSFTimeStamp(data.LastModifiedDateTime)
		isPrimary := strconv.FormatBool(data.IsPrimary)

		sqlquery := fmt.Sprintf(`IF EXISTS(SELECT personIdExternal FROM PerEmail WHERE personIdExternal = '%s' AND emailType = '%s') UPDATE PerEmail SET [emailAddress]='%s', [createdOn]='%s', [lastModifiedDateTime]='%s', [isPrimary]='%s' WHERE [personIdExternal] ='%s' AND [emailType] = '%s' ELSE INSERT INTO PerEmail (emailType,personIdExternal,emailAddress,createdOn,lastModifiedDateTime,isPrimary) VALUES('%s','%s','%s','%s','%s','%s');`, personIDExternal, emailType, emailAddress, createdOn, lastModifiedDateTime, isPrimary, personIDExternal, emailType, emailType, personIDExternal, emailAddress, createdOn, lastModifiedDateTime, isPrimary)

		db.QueryMsSql(sqlquery)
	}
	fmt.Println()
}

func SyncPerNationalID(hours int) {
	var (
		dataSlice []PerNationalID
	)
	sfApi := mapSFApi["PerNationalID"] //sfApi += `&$filter=externalCode eq '00003812'`

	if hours > 0 { // update operation
		h, _ := time.ParseDuration("-" + strconv.Itoa(hours) + "h")
		sfApi += "&$filter=lastModifiedDateTime gt datetime'" + time.Now().Add(h).Format("2006-01-02T15:04:05") + "'" // + "T00:00:00'"
	}
	sfApiNext := sfApi

	for { // query data from SF, then insert to struct slice
		jsonRaw := GetSFJsonRaw(sfApiNext)

		dataResults := new(PerNationalIDResults)          // must be renewed here, or will use the last value when meets the null value
		err := json.Unmarshal(jsonRaw["d"], &dataResults) // remove the first "d", and put result to EmployeeResult struct
		output.GenerateLog(err, "Results json parsing", string(jsonRaw["d"]), false)

		dataSlice = append(dataSlice, dataResults.Results...)

		sfApiNext = dataResults.NextURL
		if sfApiNext == "" { //here we meet the end of the query
			break
		}
	}
	output.GenerateLog(nil, "SF data retrieving", strconv.Itoa(len(dataSlice))+" records retrieved in "+sfApi, true)

	b := pgb.NewBar("PerNationalID", len(dataSlice))
	for _, data := range dataSlice {
		b.Add()

		country := data.Country
		personIDExternal := data.PersonIDExternal
		createOn := common.ConvertSFTimeStamp(data.CreatedOn)
		lastModifiedDateTime := common.ConvertSFTimeStamp(data.LastModifiedDateTime)
		customString1 := common.ConvertText(data.CustomString1)
		nationalID := data.NationalID

		sqlquery := fmt.Sprintf(`IF EXISTS(SELECT personIdExternal FROM PerNationalID WHERE personIdExternal = '%s' AND nationalId = '%s') UPDATE PerNationalID SET [country]='%s', [createdOn]='%s', [lastModifiedDateTime]='%s', [customString1]='%s' WHERE [personIdExternal] ='%s' ELSE INSERT INTO PerNationalID (country,PersonIDExternal,createdOn,lastModifiedDateTime,customString1,nationalId) VALUES('%s','%s','%s','%s','%s','%s');`, personIDExternal, nationalID, country, createOn, lastModifiedDateTime, customString1, personIDExternal, country, personIDExternal, createOn, lastModifiedDateTime, customString1, nationalID)

		db.QueryMsSql(sqlquery)
	}
	fmt.Println()
}

func SyncPerPhone(hours int) {
	var (
		dataSlice []PerPhone
	)
	sfApi := mapSFApi["PerPhone"] //sfApi += `&$filter=externalCode eq '00003812'`

	if hours > 0 { // update operation
		h, _ := time.ParseDuration("-" + strconv.Itoa(hours) + "h")
		sfApi += "&$filter=lastModifiedDateTime gt datetime'" + time.Now().Add(h).Format("2006-01-02T15:04:05") + "'" // + "T00:00:00'"
	}
	sfApiNext := sfApi

	for { // query data from SF, then insert to struct slice
		jsonRaw := GetSFJsonRaw(sfApiNext)

		dataResults := new(PerPhoneResults)               // must be renewed here, or will use the last value when meets the null value
		err := json.Unmarshal(jsonRaw["d"], &dataResults) // remove the first "d", and put result to EmployeeResult struct
		output.GenerateLog(err, "Results json parsing", string(jsonRaw["d"]), false)

		dataSlice = append(dataSlice, dataResults.Results...)

		sfApiNext = dataResults.NextURL
		if sfApiNext == "" { //here we meet the end of the query
			break
		}
	}
	output.GenerateLog(nil, "SF data retrieving", strconv.Itoa(len(dataSlice))+" records retrieved in "+sfApi, true)

	b := pgb.NewBar("PerPhone", len(dataSlice))
	for _, data := range dataSlice {
		b.Add()

		phoneType := data.PhoneType
		personIDExternal := data.PersonIDExternal
		phoneNumber := common.ConvertText(data.PhoneNumber)
		createdOn := common.ConvertSFTimeStamp(data.CreatedOn)
		lastModifiedDateTime := common.ConvertSFTimeStamp(data.LastModifiedDateTime)
		isPrimary := strconv.FormatBool(data.IsPrimary)
		countryCode := data.CountryCode

		sqlquery := fmt.Sprintf(`IF EXISTS(SELECT personIdExternal FROM PerPhone WHERE personIdExternal = '%s' AND phoneType = '%s') UPDATE PerPhone SET [phoneNumber]='%s', [createdOn]='%s', [lastModifiedDateTime]='%s', [isPrimary]='%s', [countryCode]='%s' WHERE [personIdExternal] ='%s' AND [phoneType] = '%s' ELSE INSERT INTO PerPhone (phoneType,personIdExternal,phoneNumber,createdOn,lastModifiedDateTime,isPrimary,countryCode) VALUES('%s','%s','%s','%s','%s','%s','%s');`, personIDExternal, phoneType, phoneNumber, createdOn, lastModifiedDateTime, isPrimary, countryCode, personIDExternal, phoneType, phoneType, personIDExternal, phoneNumber, createdOn, lastModifiedDateTime, isPrimary, countryCode)

		db.QueryMsSql(sqlquery)
	}
	fmt.Println()
}

func SyncPosition(hours int) {
	var (
		dataSlice []Position
	)
	sfApi := mapSFApi["Position"] //sfApi += `&$filter=externalCode eq '00003812'`

	if hours > 0 { // update operation
		h, _ := time.ParseDuration("-" + strconv.Itoa(hours) + "h")
		sfApi += "&$filter=effectiveStartDate gt datetime'" + time.Now().Add(h).Format("2006-01-02T15:04:05") + "'" // + "T00:00:00'"
	}
	sfApiNext := sfApi

	for { // query data from SF, then insert to struct slice
		jsonRaw := GetSFJsonRaw(sfApiNext)

		dataResults := new(PositionResults)               // must be renewed here, or will use the last value when meets the null value
		err := json.Unmarshal(jsonRaw["d"], &dataResults) // remove the first "d", and put result to EmployeeResult struct
		output.GenerateLog(err, "Results json parsing", string(jsonRaw["d"]), false)

		dataSlice = append(dataSlice, dataResults.Results...)

		sfApiNext = dataResults.NextURL
		if sfApiNext == "" { //here we meet the end of the query
			break
		}
	}
	output.GenerateLog(nil, "SF data retrieving", strconv.Itoa(len(dataSlice))+" records retrieved in "+sfApi, true)

	b := pgb.NewBar("Position", len(dataSlice))
	for _, data := range dataSlice {
		b.Add()

		code := data.Code
		effectiveStartDate := common.ConvertSFTimeStamp(data.EffectiveStartDate)
		jobCode := common.ConvertText(data.JobCode)
		positionType := common.ConvertText(data.Type)
		division := common.ConvertText(data.Division)
		externalNameLocalized := common.ConvertText(data.ExternalNameLocalized)
		effectiveStatus := common.ConvertText(data.EffectiveStatus)
		description := common.ConvertText(data.Description)
		externalNameDefaultValue := common.ConvertText(data.ExternalNameDefaultValue)
		positionControlled := strconv.FormatBool(data.PositionControlled)
		company := common.ConvertText(data.Company)
		department := common.ConvertText(data.Department)
		targetFTE := common.ConvertText(data.TargetFTE)
		jobLevel := common.ConvertText(data.JobLevel)
		externalNameEnUS := common.ConvertText(data.ExternalNameEnUS)

		sqlquery := fmt.Sprintf(`IF EXISTS(SELECT code FROM Position WHERE code = '%s') UPDATE Position SET [effectiveStartDate]='%s', [jobCode]='%s', [type]='%s', [division]='%s', [externalName_localized]='%s', [effectiveStatus]='%s', [description]='%s', [externalName_defaultValue]='%s', [positionControlled]='%s', [company]='%s', [department]='%s', [targetFTE]='%s', [jobLevel]='%s', [externalName_en_US]='%s' WHERE [code] ='%s' ELSE INSERT INTO Position (code, effectiveStartDate, jobCode, type, division, externalName_localized, effectiveStatus, description, externalName_defaultValue, positionControlled, company, department, targetFTE, jobLevel, externalName_en_US) VALUES('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s');`, code, effectiveStartDate, jobCode, positionType, division, externalNameLocalized, effectiveStatus, description, externalNameDefaultValue, positionControlled, company, department, targetFTE, jobLevel, externalNameEnUS, code, code, effectiveStartDate, jobCode, positionType, division, externalNameLocalized, effectiveStatus, description, externalNameDefaultValue, positionControlled, company, department, targetFTE, jobLevel, externalNameEnUS)

		db.QueryMsSql(sqlquery)
	}
	fmt.Println()
}

func SyncUser(hours int) {
	var (
		dataSlice []User
	)
	sfApi := mapSFApi["User"] //sfApi += `&$filter=externalCode eq '00003812'`

	if hours > 0 { // update operation
		h, _ := time.ParseDuration("-" + strconv.Itoa(hours) + "h")
		sfApi += "&$filter=lastModifiedDateTime gt datetime'" + time.Now().Add(h).Format("2006-01-02T15:04:05") + "'" // + "T00:00:00'"
	}
	sfApiNext := sfApi

	for { // query data from SF, then insert to struct slice
		jsonRaw := GetSFJsonRaw(sfApiNext)

		dataResults := new(UserResults)                   // must be renewed here, or will use the last value when meets the null value
		err := json.Unmarshal(jsonRaw["d"], &dataResults) // remove the first "d", and put result to EmployeeResult struct
		output.GenerateLog(err, "Results json parsing", string(jsonRaw["d"]), false)

		dataSlice = append(dataSlice, dataResults.Results...)

		sfApiNext = dataResults.NextURL
		if sfApiNext == "" { //here we meet the end of the query
			break
		}
	}
	output.GenerateLog(nil, "SF data retrieving", strconv.Itoa(len(dataSlice))+" records retrieved in "+sfApi, true)

	b := pgb.NewBar("User", len(dataSlice))
	for _, data := range dataSlice {
		b.Add()

		userID := common.ConvertText(data.UserID)
		division := common.ConvertText(data.Division)
		custom01 := common.ConvertText(data.Custom01)
		lastModifiedDateTime := common.ConvertSFTimeStamp(data.LastModifiedDateTime)
		timeZone := common.ConvertText(data.TimeZone)
		defaultLocale := common.ConvertText(data.DefaultLocale)
		status := common.ConvertText(data.Status)
		lastName := common.ConvertText(data.LastName)
		email := common.ConvertText(data.Email)
		defaultFullname := common.ConvertText(data.DefaultFullName)
		country := common.ConvertText(data.Country)
		department := common.ConvertText(data.Department)
		firstName := common.ConvertText(data.FirstName)
		empID := common.ConvertText(data.EmpID)
		title := common.ConvertText(data.Title)
		hireDate := common.ConvertSFTimeStamp(data.HireDate)
		location := common.ConvertText(data.Location)
		username := common.ConvertText(data.Username)

		sqlquery := fmt.Sprintf(`IF EXISTS(SELECT userId FROM Users WHERE userId = '%s') UPDATE Users SET [division]='%s', [custom01]='%s', [lastModifiedDateTime]='%s', [timeZone]='%s', [defaultLocale]='%s', [status]='%s', [lastName]=N'%s', [email]='%s', [defaultFullName]=N'%s', [country]='%s', [department]='%s', [firstName]=N'%s', [empID]='%s', [title]='%s' ,[hireDate]='%s', [location]='%s', [username]='%s' WHERE [userId] ='%s' ELSE INSERT INTO Users (userId, division, custom01, lastModifiedDateTime, timeZone, defaultLocale, status, lastName, email, defaultFullName, country, department, firstName, empId, title, hireDate, location, username) VALUES('%s','%s','%s','%s','%s','%s','%s',N'%s','%s',N'%s','%s','%s',N'%s','%s','%s','%s','%s','%s');`, userID, division, custom01, lastModifiedDateTime, timeZone, defaultLocale, status, lastName, email, defaultFullname, country, department, firstName, empID, title, hireDate, location, username, userID, userID, division, custom01, lastModifiedDateTime, timeZone, defaultLocale, status, lastName, email, defaultFullname, country, department, firstName, empID, title, hireDate, location, username)

		db.QueryMsSql(sqlquery)
	}
	fmt.Println()
}

func SyncServiceCompany() {
	var (
		dataSlice []ServiceCompany
	)
	sfApi := mapSFApi["ServiceCompany"] //sfApi += `&$filter=externalCode eq '00003812'`

	sfApiNext := sfApi

	for { // query data from SF, then insert to struct slice
		jsonRaw := GetSFJsonRaw(sfApiNext)

		dataResults := new(ServiceCompanyResults)     // must be renewed here, or will use the last value when meets the null value
		err := json.Unmarshal(jsonRaw["d"], &jsonRaw) // remove the first "d", and put result to EmployeeResult struct
		output.GenerateLog(err, "Results json parsing", string(jsonRaw["d"]), false)

		err = json.Unmarshal(jsonRaw["picklistOptions"], &dataResults)
		output.GenerateLog(err, "Results json parsing", string(jsonRaw["picklistOptions"]), false)

		dataSlice = append(dataSlice, dataResults.Results...)

		sfApiNext = dataResults.NextURL
		if sfApiNext == "" { //here we meet the end of the query
			break
		}
	}
	output.GenerateLog(nil, "SF data retrieving", strconv.Itoa(len(dataSlice))+" records retrieved in "+sfApi, true)

	b := pgb.NewBar("ServiceCompany", len(dataSlice))
	for _, data := range dataSlice {
		b.Add()

		id := data.ID
		externalCode := data.ExternalCode

		sqlquery := fmt.Sprintf(`IF EXISTS(SELECT id FROM ServiceCompany WHERE id = '%s') UPDATE ServiceCompany SET [externalCode]='%s' WHERE [id] ='%s' ELSE INSERT INTO ServiceCompany (id, externalCode) VALUES('%s','%s');`, id, externalCode, id, id, externalCode)

		db.QueryMsSql(sqlquery)
	}
	fmt.Println()
}
