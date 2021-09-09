package db

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"time"
)

const (
	server   = "10.253.32.52"
	port     = 3306
	user     = "rosicky.gui"
	password = "rosicky.gui"
	database = "UNI_ID"
	dbString = "server=%s;port%d;database=%s;user id=%s;password=%s;encrypt=disable;connection timeout=300;dial timeout=300;"
)

var (
	connString = fmt.Sprintf(dbString, server, port, database, user, password)
)

type SFView struct {
	PersonIDExternal string
	StartDate        string
	CreatedOn        string
	LastDateWorked   string
	Site             string
	CompanyEN        sql.NullString //maybe null in db
	CompanyCN        sql.NullString //maybe null in db
	AddressEN        sql.NullString //maybe null in db
	AddressCN        sql.NullString //maybe null in db
	EmplStatus       int
	ManagerID        string
	FirstName        sql.NullString //maybe null in db
	LastName         sql.NullString //maybe null in db
	Username         string
	EmailAddress     sql.NullString
	PhoneNumber      sql.NullString //maybe null in db
	SignLang         int
	JobLevel         sql.NullString //maybe null in db
	Department       sql.NullString //maybe null in db
	DepartmentCN     sql.NullString //maybe null in db
	TitleEN          sql.NullString //maybe null in db
	TitleCN          sql.NullString //maybe null in db
	JobEffectDate    string
	PosEffectDate    string
	LastModDate      string
}

type ADUserAttr struct {
	SAMAccountName             string `gorm:"primary_key"`
	EmployeeID                 string
	PhysicalDeliveryOfficeName string
	Pager                      string
	Sn                         string
	Mobile                     string
	Company                    string
	Title                      string
	Department                 string
	Comment                    string
	DisplayName                string
	Description                string
	Givenname                  string
	UserPrincipalName          string
	Mail                       string
	Manager                    string
	UserAccountControl         string
	MSExchHideFromAddressLists string

	LastOperation string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
