package db

import (
	"fmt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"uidm/common"
	"uidm/output"
)

func ConnectORM() *gorm.DB {
	dsn := fmt.Sprintf(dbString, server, port, database, user, password)
	orm, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	output.GenerateLog(err, "GORM connection", dsn, false)

	return orm
}

func CreateORM(orm *gorm.DB, u map[string][]string) {
	a := ADUserAttr{
		SAMAccountName:             common.GetStringSliceFirst(u["sAMAccountName"]),
		EmployeeID:                 common.GetStringSliceFirst(u["employeeID"]),
		PhysicalDeliveryOfficeName: common.GetStringSliceFirst(u["physicalDeliveryOfficeName"]),
		Mobile:                     common.GetStringSliceFirst(u["mobile"]),
		Pager:                      common.GetStringSliceFirst(u["pager"]),
		Sn:                         common.GetStringSliceFirst(u["sn"]),
		Company:                    common.GetStringSliceFirst(u["company"]),
		Title:                      common.GetStringSliceFirst(u["title"]),
		Department:                 common.GetStringSliceFirst(u["department"]),
		Comment:                    common.GetStringSliceFirst(u["comment"]),
		DisplayName:                common.GetStringSliceFirst(u["displayName"]),
		Description:                common.GetStringSliceFirst(u["description"]),
		Givenname:                  common.GetStringSliceFirst(u["givenname"]),
		UserPrincipalName:          common.GetStringSliceFirst(u["userPrincipalName"]),
		Mail:                       common.GetStringSliceFirst(u["mail"]),
		Manager:                    common.ConvertDNToCN(common.GetStringSliceFirst(u["manager"])),
		UserAccountControl:         common.GetStringSliceFirst(u["userAccountControl"]),
		MSExchHideFromAddressLists: "FALSE",

		LastOperation: "create",
	}
	orm.Create(&a)
}

func UpdateORM(orm *gorm.DB, sam string, u map[string][]string) {
	var a = new(ADUserAttr)

	orm.Where("sam_account_name = ?", sam).Find(&a)

	a.PhysicalDeliveryOfficeName = common.GetStringSliceFirst(u["physicalDeliveryOfficeName"])
	a.Mobile = common.GetStringSliceFirst(u["mobile"])
	a.Pager = common.GetStringSliceFirst(u["pager"])
	a.Sn = common.GetStringSliceFirst(u["sn"])
	a.Company = common.GetStringSliceFirst(u["company"])
	a.Title = common.GetStringSliceFirst(u["title"])
	a.Department = common.GetStringSliceFirst(u["department"])
	a.Comment = common.GetStringSliceFirst(u["comment"])
	a.DisplayName = common.GetStringSliceFirst(u["displayName"])
	a.Description = common.GetStringSliceFirst(u["description"])
	a.Givenname = common.GetStringSliceFirst(u["givenname"])
	a.Manager = common.ConvertDNToCN(common.GetStringSliceFirst(u["manager"]))

	a.LastOperation = "update"

	orm.Save(&a)
}

func DeleteORM(orm *gorm.DB, sam string, u map[string][]string) {
	var a = new(ADUserAttr)
	orm.Where("sam_account_name = ?", sam).Find(&a)

	a.UserAccountControl = common.GetStringSliceFirst(u["userAccountControl"])
	a.Manager = ""
	a.Description = common.GetStringSliceFirst(u["description"])
	a.MSExchHideFromAddressLists = "TRUE"

	a.LastOperation = "disable"

	orm.Save(&a)
	orm.Delete(&a)
}
