package sf

import "github.com/qianlnk/pgbar"

const (
	sfAuth = "Basic YXBpLml0QGNzaXNvbGFycG86Q3NpQDIwMTU="
)

var (
	mapSFApi = map[string]string{
		"PerPersonal":    "https://api15.sapsf.cn:443/odata/v2/PerPersonal?$format=JSON&$orderby=personIdExternal",
		"FOCompany":      "https://api15.sapsf.cn:443/odata/v2/FOCompany?$format=JSON&$orderby=externalCode",
		"EmpEmployment":  "https://api15.sapsf.cn:443/odata/v2/EmpEmployment?$format=JSON&$orderby=personIdExternal",
		"EmpJob":         "https://api15.sapsf.cn:443/odata/v2/EmpJob?$format=JSON&$orderby=userId",
		"FODepartment":   "https://api15.sapsf.cn:443/odata/v2/FODepartment?$format=JSON&$orderby=externalCode",
		"PerEmail":       "https://api15.sapsf.cn:443/odata/v2/PerEmail?$format=JSON&$orderby=personIdExternal",
		"PerNationalID":  "https://api15.sapsf.cn:443/odata/v2/PerNationalId?$format=JSON&$orderby=personIdExternal",
		"PerPhone":       "https://api15.sapsf.cn:443/odata/v2/PerPhone?$format=JSON&$orderby=personIdExternal",
		"Position":       "https://api15.sapsf.cn:443/odata/v2/Position?$format=JSON&$orderby=code",
		"User":           "https://api15.sapsf.cn:443/odata/v2/User?$format=JSON&$orderby=userId",
		"ServiceCompany": "https://api15.sapsf.cn:443/odata/v2/Picklist('serviceCompany')?$expand=picklistOptions&$format=JSON&$orderby=id",
	}
	pgb = pgbar.New("")
)

// PerPersonal struct
type PerPersonal struct {
	PersonIDExternal     string `json:"personIdExternal"`
	StartDate            string `json:"startDate"`
	LastName             string `json:"lastName"`
	Gender               string `json:"gender"`
	EndDate              string `json:"endDate"`
	CreatedOn            string `json:"createdOn"`
	CustomString5        string `json:"customString5"`
	CustomString6        string `json:"customString6"`
	LastModifiedDateTime string `json:"lastModifiedDateTime"`
	FirstName            string `json:"firstName"`
}

// PerPersonalResults struct
type PerPersonalResults struct {
	Results []PerPersonal `json:"results"`
	NextURL string        `json:"__next"`
}

// FOCompany struct
type FOCompany struct {
	ExternalCode string `json:"externalCode"`

	StartDate            string `json:"startDate"`
	Country              string `json:"country"`
	NameLocalized        string `json:"name_localized"`
	NameThTH             string `json:"name_th_TH"`
	Name                 string `json:"name"`
	Status               string `json:"status"`
	Description          string `json:"description"`
	EndDate              string `json:"endDate"`
	CreatedOn            string `json:"createdOn"`
	LastModifiedDateTime string `json:"lastModifiedDateTime"`
	NameEnUS             string `json:"name_en_US"`
	Currency             string `json:"currency"`
	NameZhCN             string `json:"name_zh_CN"`
	DefaultLocation      string `json:"defaultLocation"`
	NameDefaultValue     string `json:"name_defaultValue"`
}

// FOCompanyResults struct
type FOCompanyResults struct {
	Results []FOCompany `json:"results"`
	NextURL string      `json:"__next"`
}

// EmpEmployment struct
type EmpEmployment struct {
	PersonIDExternal     string `json:"personIdExternal"`
	UserID               string `json:"userId"`
	EndDate              string `json:"endDate"`
	AssignmentIDExternal string `json:"assignmentIdExternal"`
	SeniorityDate        string `json:"seniorityDate"`
	StartDate            string `json:"startDate"`
	HiringNotCompleted   bool   `json:"hiringNotCompleted"`
	IsECRecord           bool   `json:"isECRecord"`
	CustomString20       string `json:"customString20"`
	CreatedOn            string `json:"createdOn"`
	LastDateWorked       string `json:"lastDateWorked"`
	OriginalStartDate    string `json:"originalStartDate"`
	AssignmentClass      string `json:"assignmentClass"`
	CustomString5        string `json:"customString5"`
	LastModifiedDateTime string `json:"lastModifiedDateTime"`
}

// EmpEmploymentResults struct
type EmpEmploymentResults struct {
	Results []EmpEmployment `json:"results"`
	NextURL string          `json:"__next"`
}

// EmpJob struct
type EmpJob struct {
	UserID               string `json:"userId"`
	StartDate            string `json:"startDate"`
	JobCode              string `json:"jobCode"`
	Division             string `json:"division"`
	EmplStatus           string `json:"emplStatus"`
	CountryOfCompany     string `json:"countryOfCompany"`
	ManagerID            string `json:"managerId"`
	LastModifiedDateTime string `json:"lastModifiedDateTime"`
	CreatedOn            string `json:"createdOn"`
	Position             string `json:"position"`
	Company              string `json:"company"`
	Department           string `json:"department"`
	EmployeeClass        string `json:"employeeClass"`
	Location             string `json:"location"`
	JobTitle             string `json:"jobTitle"`
	CustomString10       string `json:"customString10"`
}

// EmpJobResults struct
type EmpJobResults struct {
	Results []EmpJob `json:"results"`
	NextURL string   `json:"__next"`
}

// FODepartment struct
type FODepartment struct {
	ExternalCode            string `json:"externalCode"`
	StartDate               string `json:"startDate"`
	Parent                  string `json:"parent"`
	NameLocalized           string `json:"name_localized"`
	CustDeptLevel           string `json:"cust_deptLevel"`
	DescriptionDefaultValue string `json:"description_defaultValue"`
	Name                    string `json:"name"`
	DescriptionEnUS         string `json:"description_en_US"`
	Status                  string `json:"status"`
	Description             string `json:"description"`
	CreatedOn               string `json:"createdOn"`
	NameEnUS                string `json:"name_en_US"`
	LastModifiedDateTime    string `json:"lastModifiedDateTime"`
	DescriptionLocalized    string `json:"description_localized"`
	CustHeadPosition        string `json:"cust_headPosition"`
	HeadOfUnit              string `json:"headOfUnit"`
}

// FODepartmentResults struct
type FODepartmentResults struct {
	Results []FODepartment `json:"results"`
	NextURL string         `json:"__next"`
}

// PerEmail struct
type PerEmail struct {
	EmailType            string `json:"emailType"`
	PersonIDExternal     string `json:"personIdExternal"`
	EmailAddress         string `json:"emailAddress"`
	CreatedOn            string `json:"createdOn"`
	LastModifiedDateTime string `json:"lastModifiedDateTime"`
	IsPrimary            bool   `json:"isPrimary"`
}

// PerMailResults struct
type PerEmailResults struct {
	Results []PerEmail `json:"results"`
	NextURL string     `json:"__next"`
}

// PerNationalID struct
type PerNationalID struct {
	Country              string `json:"country"`
	PersonIDExternal     string `json:"personIdExternal"`
	CreatedOn            string `json:"createdOn"`
	LastModifiedDateTime string `json:"lastModifiedDateTime"`
	CustomString1        string `json:"customString1"`
	NationalID           string `json:"nationalId"`
}

// PerNationalIDResults struct
type PerNationalIDResults struct {
	Results []PerNationalID `json:"results"`
	NextURL string          `json:"__next"`
}

// PerPhone struct
type PerPhone struct {
	PhoneType            string `json:"phoneType"`
	PersonIDExternal     string `json:"personIdExternal"`
	PhoneNumber          string `json:"phoneNumber"`
	CreatedOn            string `json:"createdOn"`
	LastModifiedDateTime string `json:"lastModifiedDateTime"`
	IsPrimary            bool   `json:"isPrimary"`
	CountryCode          string `json:"countryCode"`
}

// PerPhoneResults struct
type PerPhoneResults struct {
	Results []PerPhone `json:"results"`
	NextURL string     `json:"__next"`
}

// Position struct
type Position struct {
	Code                     string `json:"code"`
	EffectiveStartDate       string `json:"effectiveStartDate"`
	JobCode                  string `json:"jobCode"`
	Type                     string `json:"type"`
	Division                 string `json:"division"`
	ExternalNameLocalized    string `json:"externalName_localized"`
	EffectiveStatus          string `json:"effectiveStatus"`
	Description              string `json:"description"`
	ExternalNameDefaultValue string `json:"externalName_defaultValue"`
	PositionControlled       bool   `json:"positionControlled"`
	Company                  string `json:"company"`
	Department               string `json:"department"`
	TargetFTE                string `json:"targetFTE"`
	JobLevel                 string `json:"jobLevel"`
	ExternalNameEnUS         string `json:"externalName_en_US"`
}

// PositionResults struct
type PositionResults struct {
	Results []Position `json:"results"`
	NextURL string     `json:"__next"`
}

// User struct
type User struct {
	UserID               string `json:"userId"`
	Division             string `json:"division"`
	Custom01             string `json:"custom01"`
	LastModifiedDateTime string `json:"lastModifiedDateTime"`
	TimeZone             string `json:"timeZone"`
	DefaultLocale        string `json:"defaultLocale"`
	Status               string `json:"status"`
	LastName             string `json:"lastName"`
	Email                string `json:"email"`
	DefaultFullName      string `json:"defaultFullName"`
	Country              string `json:"country"`
	Department           string `json:"department"`
	FirstName            string `json:"firstName"`
	EmpID                string `json:"empId"`
	Title                string `json:"title"`
	HireDate             string `json:"hireDate"`
	Location             string `json:"location"`
	Username             string `json:"username"`
}

// UserResults struct
type UserResults struct {
	Results []User `json:"results"`
	NextURL string `json:"__next"`
}

// ServiceCompany struct
type ServiceCompany struct {
	ID           string `json:"id"`
	ExternalCode string `json:"externalCode"`
}

// ServiceCompanyResults struct
type ServiceCompanyResults struct {
	Results []ServiceCompany `json:"results"`
	NextURL string           `json:"__next"`
}
