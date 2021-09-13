package ad

import "github.com/qianlnk/pgbar"

const (
	disGrpOU = "OU=Distribution Groups,OU=Groups,OU=Control"
	secGrpOU = "OU=Security Groups,OU=Groups,OU=Control,OU=SZSP,OU=Suzhou,OU=China,OU=APAC"

	defaultPWD = "Csi@2021"

	siteFilter = ` AND site in ('SZSP', 'SZST', 'CSSC', 'SZGD', 'SZSM', 'SZEM', 'SZED', 'SZES', 'SZIV', 'YCSM', 'YCPT', 'YCDF', 'SZCC', 'SZSE','YCSE','LYSP','LYPT', 'SQSE', 'CSAS', 'BTSE', 'JXSE', 'JXST', 'JXNM', 'CSSM', 'CSTG', 'CSTL', 'XNPT','SSCH')`
	//siteFilter = ` AND username = 'Y1E03710'`
	//siteFilter = ` AND site in ('YCSE')`
	personFilter = ` AND (titleCn is NULL or titleCn not like '%总经理') AND (emailaddress NOT LIKE '%@canadiansolar.com' OR emailAddress is NULL) AND personIdExternal not in ('S0E01763', 'S1E00005', 'C1E14283', 'S0E00001', 'S0E00002', 'S0E01147', 'S0E00626','S0E00919')`

	queryUpdateAll  = `SELECT * from UNI_ID..View_SFAttrForAD WHERE username != 'null' AND lastDateWorked < startDate`
	queryDisableAll = `SELECT * FROM UNI_ID..View_SFAttrForAD WHERE username != 'null' AND lastDateWorked >= startDate AND lastDateWorked <= GETDATE()`
)

var (
	pgb = pgbar.New("")

	MapLDAP = map[string]map[string]string{
		"dev": {
			"url":    "ldaps://10.253.32.198:636",
			"adm":    "CN=Administrator,CN=Users",
			"pwd":    "Csi@solar10",
			"baseDN": "DC=uidtest,DC=local",
		},
		"prd": {
			"url":    "ldaps://10.0.5.5:636",
			"adm":    "CN=service.monitor,OU=Service Accounts,OU=Control,OU=SZSP,OU=Suzhou,OU=China,OU=APAC",
			"pwd":    "Csi@solar10",
			"baseDN": "DC=csisolar,DC=com",
		},
	}

	mapSiteGroup = map[string]map[string]string{
		"BTSE": {
			"group": "CN_BTSE_ALL",
			"ou":    "OU=Baotou,OU=China,OU=APAC",
		},
		"CSAS": {
			"group": "CN_CSAS_ALL",
			"ou":    "OU=Changshu,OU=China,OU=APAC",
		},
		"CSSM": {
			"group": "CN_CSSM_ALL",
			"ou":    "OU=Changshu,OU=China,OU=APAC",
		},
		"CSSC": {
			"group": "CN_CSSC_ALL",
			"ou":    "OU=Suzhou,OU=China,OU=APAC",
		},
		"CSTG": {
			"group": "CN_CSTG_ALL",
			"ou":    "OU=Changshu,OU=China,OU=APAC",
		},
		"CSTL": {
			"group": "CN_CSTL_ALL",
			"ou":    "OU=Changshu,OU=China,OU=APAC",
		},
		"JXSE": {
			"group": "CN_JXSE_ALL",
			"ou":    "OU=Jiaxing,OU=China,OU=APAC",
		},
		"JXST": {
			"group": "CN_JXST_ALL",
			"ou":    "OU=Jiaxing,OU=China,OU=APAC",
		},
		"JXNM": {
			"group": "CN_JXNM_ALL",
			"ou":    "OU=Jiaxing,OU=China,OU=APAC",
		},
		"LYPT": {
			"group": "CN_LYPT_ALL",
			"ou":    "OU=Luoyang,OU=China,OU=APAC",
		},
		"LYSP": {
			"group": "CN_LYSP_ALL",
			"ou":    "OU=Luoyang,OU=China,OU=APAC",
		},
		"SQSE": {
			"group": "CN_SQSE_ALL",
			"ou":    "OU=Suqian,OU=China,OU=APAC",
		},
		"SZCC": {
			"group": "CN_SZCC_ALL",
			"ou":    "OU=Suzhou,OU=China,OU=APAC",
		},
		"SZED": {
			"group": "CN_SZED_ALL",
			"ou":    "OU=Suzhou,OU=China,OU=APAC",
		},
		"SZES": {
			"group": "CN_SZES_ALL",
			"ou":    "OU=Suzhou,OU=China,OU=APAC",
		},
		"SZET": {
			"group": "CN_SZET_ALL",
			"ou":    "OU=Suzhou,OU=China,OU=APAC",
		},
		"SZIV": {
			"group": "CN_SZIV_ALL",
			"ou":    "OU=Suzhou,OU=China,OU=APAC",
		},
		"SZSE": {
			"group": "CN_SZSE_ALL",
			"ou":    "OU=Suzhou,OU=China,OU=APAC",
		},
		"SZGD": {
			"group": "CN_SZGD_ALL",
			"ou":    "OU=Suzhou,OU=China,OU=APAC",
		},
		"SZSM": {
			"group": "CN_SZSM_ALL",
			"ou":    "OU=Suzhou,OU=China,OU=APAC",
		},
		"SZSP": {
			"group": "CN_SZSP_ALL",
			"ou":    "OU=Suzhou,OU=China,OU=APAC",
		},
		"SZST": {
			"group": "CN_SZST_ALL",
			"ou":    "OU=Suzhou,OU=China,OU=APAC",
		},
		"SSCH": {
			"group": "CN_SSCH_ALL",
			"ou":    "OU=Suzhou,OU=China,OU=APAC",
		},
		"THSM": {
			"group": "CN_THSM_ALL",
			"ou":    "OU=Chonburi,OU=Thailand,OU=APAC",
		},
		"VNSM": {
			"group": "CN_VNSM_ALL",
			"ou":    "OU=HaiPhong,OU=Vietnam,OU=APAC",
		},
		"YCDF": {
			"group": "CN_YCDF_ALL",
			"ou":    "OU=Yancheng,OU=China,OU=APAC",
		},
		"YCSE": {
			"group": "CN_YCSE_ALL",
			"ou":    "OU=Yancheng,OU=China,OU=APAC",
		},
		"YCSM": {
			"group": "CN_YCSM_ALL",
			"ou":    "OU=Yancheng,OU=China,OU=APAC",
		},
		"YCPT": {
			"group": "CN_YCPT_ALL",
			"ou":    "OU=Yancheng,OU=China,OU=APAC",
		},
		"HKCS": {
			"group": "CN_HKCS_ALL",
			"ou":    "OU=Hongkong,OU=China,OU=APAC",
		},
		"HKMS": {
			"group": "CN_HKMS_ALL",
			"ou":    "OU=Hongkong,OU=China,OU=APAC",
		},
		"HKEH": {
			"group": "CN_HKEH_ALL",
			"ou":    "OU=Honkgong,OU=China,OU=APAC",
		},
		"HKNE": {
			"group": "CN_HKNE_ALL",
			"ou":    "OU=Hongkong,OU=China,OU=APAC",
		},
		"XNPT": {
			"group": "CN_XNPT_ALL",
			"ou":    "OU=Xining,OU=China,OU=APAC",
		},
	}

	mapExclaimer = map[string]string{
		"CN": "CN=F_Exclaimer_ChinaCN_Local",
		"EN": "CN=F_Exclaimer_ChinaEN_Local",
	}

	mapO365 = map[string]string{
		"E3":  "CN=F_O365_E3",
		"P1":  "CN=F_O365_P1",
		"MFA": "CN=F_O365_MFA",
	}

	mapZohoURL = map[string]string{
		"dev": "https://sdpfos.canadiansolar.com/api/v3/requests",
		"prd": "https://help.csisolar.com/api/v3/requests",
	}
	mapZohoKey = map[string]string{
		"dev": "C5E91ACA-5790-4A70-B3D0-264F7DA8671F",
		"prd": "934A57D3-1433-4C9A-903E-EB11C5588464",
	}

	//mapZohoMail = map[string]map[string]string{
	//	"new": {
	//		"subject": "New employee on boarding. Please setup PC - %s",
	//		"body":    "New employee on boarding, please find details as below:\nName: %s\nUser ID: %s\nSite: %s\nJob level: %s\nStart date: %s",
	//	},
	//	"update": {
	//		"subject": "New employee on boarding. Please setup PC - %s",
	//		"body":    "New employee on boarding, please find details as below:\nName: %s\nUser ID: %s\nSite: %s\nJob level: %s\nStart date: %s",
	//	},
	//}

	mapZohoData = map[string]map[string]string{
		"dev": {
			"new": `{"request":{"requester":{"name":"CSIMAIL"},"subject": "New employee on boarding. Please setup PC - %s","description":"<p>There is an employee on boarding,detail as below: Please setup PC for him/her:<br/><br/>Name: %s<br/>User ID: %s<br/>Site: %s<br/>Job level: %s<br/>Start date: %s<p/>","service_category":{"name":"Email"}}}`,

			"update": `{"request":{"requester":{"name":"CSIMAIL"},"subject": "Employee changed work site. Please confirm if PC need transfer - %s","description":"<p>There is an employee whose work site has been changed. Please find the details as below and confirm if his/her PC need to be transferred:<br/><br/>Name: %s<br/>User ID: %s<br/>Previous site: %s<br/>Current site: %s<br/>Job level: %s<br/>Job effective date: %s<p/>","service_category":{"name":"Email"}}}`,

			"disable": `{"request":{"requester":{"name":"CSIMAIL"},"subject": "Employee resigned. Please delete the canadiansolar mailbox - %s","description":"<p>There is an employee resigned whose start date is before 2020-10-01:<br/><br/>Name: %s<br/>User ID: %s<br/>Site: %s<br/>Job level: %s<br/>Start date: %s<br/>Last date worked: %s<p/>","service_category":{"name":"Email"}}}`,

			"usernamechange": `{"request":{"requester":{"name":"CSIMAIL"},"subject":"Username may be changed in SF. Please check the new AD user - %s","description":"<p>Username may be changed in SF. <br/><br/>Current username in SF: %s<br/>Old username in AD: %s<br/>EmployeeID: %s<br/>Site: %s<br/>Job level: %s<br/>Start date: %s<p/>","service_category":{"name":"Email"}}}`,
		},
		"prd": {
			"new": `{"request":{"level":{"name":"02. Site Company / 区域公司"},"group":{"name":"LV1. V_Remote Service"},"requester":{"name":"UID"},"request_type":{"name":"Request/请求"},"category":{"name":"01. Service Delivery"},"service_category":{"name":"01. Service Delivery"},"subject":"New employee on boarding. Please setup PC - %s","subcategory":{"name":"AD/Email Account"},"mode":{"name":"Web Form"},"description":"<p>There is an employee on boarding,detail as below. Please setup PC for him/her:<br/><br/>Name: %s<br/> User ID: %s<br/> Site: %s<br/>Job level: %s<br/>Start date: %s<p/>"}}`,

			"update": `{"request":{"level":{"name":"02. Site Company / 区域公司"},"group":{"name":"LV1. V_Remote Service"},"requester":{"name":"UID"},"request_type":{"name":"Request/请求"},"category":{"name":"01. Service Delivery"},"service_category":{"name":"01. Service Delivery"},"subject":"Employee changed work site. Please confirm if PC need transfer - %s","subcategory":{"name":"AD/Email Account"},"mode":{"name":"Web Form"},"description":"<p>There is an employee whose work site has been changed. Please find the details as below and confirm if his/her PC need to be transferred:<br/><br/>Name: %s<br/> UserID: %s <br/> Previous site: %s <br/> Current site: %s <br/>Job level: %s<br/>Job effective date: %s</p>"}}`,

			"disable": `{"request":{"level":{"name":"02. Site Company / 区域公司"},"group":{"name":"LV1. V_Remote Service"},"requester":{"name":"UID"},"request_type":{"name":"Request/请求"},"category":{"name":"01. Service Delivery"},"service_category":{"name":"01. Service Delivery"},"subject":"Employee resigned. Please delete the canadiansolar mailbox - %s","subcategory":{"name":"AD/Email Account"},"mode":{"name":"Web Form"},"description":"<p>There is an employee resigned whose start date is before 2020-10-01. <br/><br/>Name: %s<br/>User ID: %s<br/>Site: %s<br/>Job level: %s<br/>Start date: %s<br/>Last date worked: %s<p/>"}}`,

			"usernamechange": `{"request":{"level":{"name":"02. Site Company / 区域公司"},"group":{"name":"LV1. V_Remote Service"},"requester":{"name":"UID"},"request_type":{"name":"Request/请求"},"category":{"name":"01. Service Delivery"},"service_category":{"name":"01. Service Delivery"},"subject":"Username may be changed in SF. Please check the new AD user - %s","subcategory":{"name":"AD/Email Account"},"mode":{"name":"Web Form"},"description":"<p>Username may be changed in SF. <br/><br/>Current username in SF: %s<br/>Old username in AD: %s<br/>EmployeeID: %s<br/>Site: %s<br/>Job level: %s<br/>Start date: %s<p/>"}}`,
		},
	}
)

type ADUser struct {
	objectClass                []string
	sAMAccountName             []string
	userAccountControl         []string
	instanceType               []string
	manager                    []string
	comment                    []string
	company                    []string
	title                      []string
	department                 []string
	userPrincipalName          []string
	mail                       []string
	employeeID                 []string
	physicalDeliveryOfficeName []string
	mobile                     []string
	pager                      []string
	sn                         []string
	displayName                []string
	description                []string
	givenname                  []string
	accountExpires             []string
	unicodePwd                 []string
	pwdLastSet                 []string
}
