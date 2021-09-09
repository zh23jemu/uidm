package ad

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/go-ldap/ldap/v3"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"strings"
	"time"
	"uidm/common"
	"uidm/db"
	"uidm/output"
)

// To generate an AD user map filled with AD user attributes
func GenerateADUserMap(l *ldap.Conn, s *db.SFView, typ, env string) map[string][]string {
	u := make(map[string][]string)

	if typ == "new" || typ == "update" {
		u["physicalDeliveryOfficeName"] = []string{s.Site}
		u["mobile"] = []string{common.ConvertMobile(s.PhoneNumber.String)}
		u["pager"] = []string{s.PersonIDExternal}
		u["sn"] = []string{s.LastName.String}
		u["userPrincipalName"] = []string{s.Username + "@csisolar.com"}

		// Signature language
		slang := "CN"
		if s.SignLang == 857 {
			if s.CompanyCN.Valid && strings.TrimSpace(s.CompanyCN.String) != "" {
				u["company"] = []string{s.CompanyCN.String}
			}
			if s.TitleCN.Valid && strings.TrimSpace(s.TitleCN.String) != "" {
				u["title"] = []string{s.TitleCN.String}
			}
			if s.DepartmentCN.Valid && strings.TrimSpace(s.DepartmentCN.String) != "" {
				u["department"] = []string{s.DepartmentCN.String}
			}
		} else {
			slang = "EN"
			if s.CompanyEN.Valid && strings.TrimSpace(s.CompanyEN.String) != "" {
				u["company"] = []string{s.CompanyEN.String}
			}
			if s.TitleEN.Valid && strings.TrimSpace(s.TitleEN.String) != "" {
				titleEN, departmentEN := common.ConvertTitleEN(s.TitleEN.String)
				u["title"] = []string{titleEN}
				if departmentEN != "" { // department may be null
					u["department"] = []string{departmentEN}
				} else {
					usrEntry := SearchADObject(l, MapLDAP[env]["baseDN"], s.Username, "sAMAccountName", []string{"department"})
					//fmt.Println(usrEntry[0].GetAttributeValue("department"))
					//当现在的部门为空时，给map的department键一个空值，用来在modifyAdUser里删掉department这个属性
					if usrEntry[0].GetAttributeValue("department") != "" {
						u["department"] = []string{""}
					}
				}
			}
		}
		u["comment"] = []string{slang}

		// If lastname is Chinese character
		isCN, _ := regexp.MatchString("[\u4e00-\u9fa5]", s.LastName.String)
		if isCN {
			names := common.SplitFirstname(s.FirstName.String)
			if len(names) == 3 { //带英文名
				u["displayName"] = []string{names[0] + " " + names[1]}
				u["givenname"] = []string{strings.ToLower(names[2] + "." + names[1])}
			} else if len(names) == 2 { //不带英文名
				u["displayName"] = []string{names[0] + " " + common.ConvertLastname(names[1])}
				u["givenname"] = []string{strings.ToLower(names[0] + "." + names[1])}
			}
			//u["displayName"] = []string{common.ConvertDisplayName(s.FirstName.String)} //需修改
			u["description"] = []string{s.LastName.String}
		} else {
			u["displayName"] = []string{s.FirstName.String + " " + s.LastName.String}
			u["givenname"] = []string{s.FirstName.String}
			u["description"] = []string{s.FirstName.String + " " + s.LastName.String}
		}

		// Mail address
		if s.EmailAddress.Valid && strings.HasSuffix(s.EmailAddress.String, "csisolar.com") && strings.Contains(s.Username, ".") {
			u["mail"] = []string{s.EmailAddress.String}
		}
		// Get manager
		mgr := SearchADObject(l, MapLDAP[env]["baseDN"], s.ManagerID, "employeeID", []string{"sAMAccountName"})
		if s.ManagerID != "" && len(mgr) > 0 {
			u["manager"] = []string{mgr[0].DN}
			//u["managerSAM"] = []string{mgr[0].GetAttributeValue("sAMAccountName")}
		}

		if typ == "new" {
			u["objectClass"] = []string{"top", "organizationalPerson", "user", "person"}
			u["userAccountControl"] = []string{fmt.Sprintf("%d", 0x0200)}
			u["instanceType"] = []string{fmt.Sprintf("%d", 0x00000004)}
			u["unicodePwd"] = []string{fmt.Sprintf(GenerateADPassword(defaultPWD))}
			u["accountExpires"] = []string{fmt.Sprintf("%d", 0x00000000)}
			u["pwdLastSet"] = []string{fmt.Sprintf("%d", 0x00000000)}

			u["sAMAccountName"] = []string{strings.ToLower(s.Username)}
			u["employeeID"] = []string{s.PersonIDExternal}
		}

	} else if typ == "disable" {
		u["userAccountControl"] = []string{fmt.Sprintf("%d", 0x0202)}
		usrEntry := SearchADObject(l, MapLDAP[env]["baseDN"], s.Username, "sAMAccountName", []string{"manager"})
		if len(usrEntry) > 0 {
			//AD账户原来有经理的才清除经理
			if usrEntry[0].GetAttributeValue("manager") != "" {
				u["manager"] = []string{}
			}
		}

		var descrp string
		entryByDescrp := SearchADObject(l, MapLDAP[env]["baseDN"], s.Username, "sAMAccountName", []string{"description"})
		if len(entryByDescrp) > 0 {
			descrp = entryByDescrp[0].GetAttributeValue("description")
		}
		u["description"] = []string{descrp + " Resigned Date: " + common.ConvertDatetimeToDate(s.LastDateWorked)}

		if env == "prd" {
			u["msExchHideFromAddressLists"] = []string{"TRUE"}
		}
	}
	return u
}

// AD user update & disabling logic
func SyncADUser(l *ldap.Conn, s *db.SFView, typ, env string, orm *gorm.DB) {
	baseDN := MapLDAP[env]["baseDN"]
	usrEntry := SearchADObject(l, baseDN, s.Username, "sAMAccountName", []string{"sAMAccountName", "physicalDeliveryOfficeName", "comment", "department", "memberOf", "description", "userAccountControl", "mail", "userPrincipalName"})

	if typ == "update" {
		//tmp := GenerateADUserMap(l, s, "new", env)
		//db.CreateORM(orm, tmp)
		if len(usrEntry) == 0 { // If user does not exist, then create it
			if !s.EmailAddress.Valid || strings.ToLower(s.Username) == strings.ToLower(strings.Split(s.EmailAddress.String, "@")[0]) || !strings.HasSuffix(s.EmailAddress.String, "csisolar.com") { // mail = null or usrname = mail or email not like csisolar.com,then create ad user
				u := GenerateADUserMap(l, s, "new", env)

				usrEntryByEmpID := SearchADObject(l, baseDN, s.PersonIDExternal, "employeeID", []string{"sAMAccountName", "userAccountControl"})
				if len(usrEntryByEmpID) > 0 {
					//存在相同工号的用户，可能是SF改名导致，触发工单提醒服务台
					sam := usrEntryByEmpID[0].GetAttributeValue("sAMAccountName")
					if s.Username != sam {
						subj := "Username may be changed in SF - " + s.Username
						body := fmt.Sprintf("Username may be changed in SF. Please check the new AD user.\n\nCurrent username in SF: %s\nOld username in AD: %s\nEmployeeID: %s\nSite: %s\nJob level: %s\nStart date:%s", s.Username, sam, s.PersonIDExternal, s.Site, s.JobLevel.String, common.ConvertDatetimeToDate(s.StartDate))
						output.SendMail(output.MailReport, subj, body, "plain", "", output.MailAdmin, output.MailHelp, output.MailOAOp)
						//data := fmt.Sprintf(mapZohoData[env]["usernamechange"], s.Username, s.Username, sam, s.PersonIDExternal, s.Site, s.JobLevel.String, common.ConvertDatetimeToDate(s.StartDate))
						//output.CreateTicket(mapZohoURL[env], mapZohoKey[env], data)
					}
				}

				// Create user, then add to groups
				usrDN := fmt.Sprintf("CN=%s,OU=Standard Users,OU=Users,OU=Resources,OU=%s,%s,%s", common.ConvertFullname(s.Username), s.Site, mapSiteGroup[s.Site]["ou"], baseDN)
				if err := CreateADUser(l, u, usrDN); err == nil {
					db.CreateORM(orm, u)

					timeStartDate, _ := time.Parse("2006-01-02", common.ConvertDatetimeToDate(s.StartDate))
					if strings.Contains(s.Username, ".") && time.Now().AddDate(0, 0, -5).Before(timeStartDate) { //people started work long ago will not create ticket
						//data := fmt.Sprintf(mapZohoData[env]["new"], s.Username, s.Username, s.PersonIDExternal, s.Site, s.JobLevel.String, common.ConvertDatetimeToDate(s.StartDate))
						//output.CreateTicket(mapZohoURL[env], mapZohoKey[env], data)
						subj := fmt.Sprintf("New employee on boarding. Please setup PC - %s", s.Username)
						body := fmt.Sprintf("New employee on boarding, please find details as below:\n\nName: %s\nUser ID: %s\nSite: %s\nJob level: %s\nStart date: %s", s.Username, s.PersonIDExternal, s.Site, s.JobLevel.String, common.ConvertDatetimeToDate(s.StartDate))
						output.SendMail(output.MailReport, subj, body, "plain", "", output.MailAdmin, output.MailHelp)
					}
					siteGrpDN := fmt.Sprintf("CN=%s,%s,OU=%s,%s,%s", mapSiteGroup[s.Site]["group"], disGrpOU, s.Site, mapSiteGroup[s.Site]["ou"], baseDN)
					sLangGrpDN := mapExclaimer[u["comment"][0]] + "," + secGrpOU + "," + baseDN

					AddADGroupMember(l, baseDN, siteGrpDN, usrDN)
					AddADGroupMember(l, baseDN, sLangGrpDN, usrDN)

					if s.EmailAddress.Valid && strings.HasSuffix(s.EmailAddress.String, "csisolar.com") && strings.Contains(strings.Split(s.EmailAddress.String, "@")[0], ".") {
						o365Grp := mapO365["P1"]
						if s.JobLevel.Valid && s.JobLevel.String <= "E" {
							o365Grp = mapO365["E3"]
						}
						o365GrpDN := o365Grp + "," + secGrpOU + "," + baseDN
						mfaGrpDN := mapO365["MFA"] + "," + secGrpOU + "," + baseDN

						AddADGroupMember(l, baseDN, o365GrpDN, usrDN)
						AddADGroupMember(l, baseDN, mfaGrpDN, usrDN)
					}
				}
			} else {
				err := errors.New("Username distinct from mail prefix")
				output.GenerateLog(err, "AD user creation", s.Site+" | "+s.Username+" | "+s.EmailAddress.String, false)
			}
		} else { // If user exists, then update it
			oldUPN := usrEntry[0].GetAttributeValue("userPrincipalName")
			if !strings.Contains(oldUPN, "canadiansolar.com") { //not deal with EG user
				u := GenerateADUserMap(l, s, "update", env)

				usrOldDN := usrEntry[0].DN
				sam := s.Username

				newDept := ""
				if len(u["department"]) > 0 {
					newDept = u["department"][0]
				}

				newSite := u["physicalDeliveryOfficeName"][0]
				newSLang := u["comment"][0] // slang is modified in GenerateADUserMap function

				oldDept := usrEntry[0].GetAttributeValue("department")
				oldSite := usrEntry[0].GetAttributeValue("physicalDeliveryOfficeName")
				oldSLang := usrEntry[0].GetAttributeValue("comment")

				oldMail := usrEntry[0].GetAttributeValue("mail")

				if strings.TrimSpace(oldMail) == "" && s.EmailAddress.Valid && strings.HasSuffix(s.EmailAddress.String, "csisolar.com") && strings.Contains(strings.Split(s.EmailAddress.String, "@")[0], ".") {
					o365Grp := mapO365["P1"]
					if s.JobLevel.Valid && s.JobLevel.String <= "E" {
						o365Grp = mapO365["E3"]
					}
					o365GrpDN := o365Grp + "," + secGrpOU + "," + baseDN
					mfaGrpDN := mapO365["MFA"] + "," + secGrpOU + "," + baseDN

					AddADGroupMember(l, baseDN, o365GrpDN, usrOldDN)
					AddADGroupMember(l, baseDN, mfaGrpDN, usrOldDN)
				}

				if err := ModifyADUser(l, u, usrOldDN); err == nil {
					db.UpdateORM(orm, sam, u)
					// Adjust sign language and site groups, move user OU
					if oldSLang != newSLang {
						newSLangGrpDN := mapExclaimer[newSLang] + "," + secGrpOU + "," + baseDN
						oldSLangGrpDN := mapExclaimer[oldSLang] + "," + secGrpOU + "," + baseDN

						AddADGroupMember(l, baseDN, newSLangGrpDN, usrOldDN)
						RemoveADGroupMember(l, baseDN, oldSLangGrpDN, usrOldDN)
					} else if strings.TrimSpace(oldDept) != "" && strings.TrimSpace(newDept) != "" && oldDept != newDept {
						if strings.Contains(s.Username, ".") {
							InformADUserGroupManagers(l, baseDN, s.Username, oldDept, newDept, "dept")
						}
					}
					if oldSite != newSite {
						//newSiteGrpDN := fmt.Sprintf("CN=CN_%s_ALL,%s,OU=%s,%s,%s", newSite, disGrpOU, newSite, mapSiteGroup[newSite]["ou"], baseDN)
						newSiteGrpDN := fmt.Sprintf("CN=%s,%s,OU=%s,%s,%s", mapSiteGroup[newSite]["group"], disGrpOU, newSite, mapSiteGroup[newSite]["ou"], baseDN)
						//oldSiteGrpDN := fmt.Sprintf("CN=CN_%s_ALL,%s,OU=%s,%s,%s", oldSite, disGrpOU, oldSite, mapSiteGroup[oldSite]["ou"], baseDN)
						oldSiteGrpDN := fmt.Sprintf("CN=%s,%s,OU=%s,%s,%s", mapSiteGroup[oldSite]["group"], disGrpOU, oldSite, mapSiteGroup[oldSite]["ou"], baseDN)

						AddADGroupMember(l, baseDN, newSiteGrpDN, usrOldDN)
						RemoveADGroupMember(l, baseDN, oldSiteGrpDN, usrOldDN)

						usrNewDN := fmt.Sprintf("CN=%s,OU=Standard Users,OU=Users,OU=Resources,OU=%s,%s,%s", s.Username, newSite, mapSiteGroup[newSite]["ou"], baseDN)
						MoveADObject(l, usrOldDN, usrNewDN)

						if strings.Contains(s.Username, ".") {
							data := fmt.Sprintf(mapZohoData[env]["update"], s.Username, s.Username, s.PersonIDExternal, oldSite, newSite, s.JobLevel.String, common.ConvertDatetimeToDate(s.JobEffectDate))
							output.CreateTicket(mapZohoURL[env], mapZohoKey[env], data)

							InformADUserGroupManagers(l, baseDN, s.Username, oldSite, newSite, "site")
						}
					}
				}
			}
		}
	} else if typ == "disable" {
		//fmt.Println("hello", s.PersonIDExternal)
		usrEntryByEmpID := SearchADObject(l, baseDN, s.PersonIDExternal, "employeeID", []string{"sAMAccountName", "userAccountControl", "description", "memberOf", "userPrincipalName"})
		if len(usrEntryByEmpID) > 0 {
			//fmt.Println(s.PersonIDExternal)
			for _, e := range usrEntryByEmpID {
				oldUPN := usrEntryByEmpID[0].GetAttributeValue("userPrincipalName")
				if !strings.Contains(oldUPN, "canadiansolar.com") { //not deal with EG user
					uacCode, _ := strconv.Atoi(e.GetAttributeValue("userAccountControl"))
					usrOldDN := e.DN

					if uacCode != 514 && strings.Contains(usrOldDN, "OU=Standard Users") { // not disabled yet, in standard user OU
						u := GenerateADUserMap(l, s, "disable", env)
						if err := DisableADUser(l, usrOldDN); err == nil {
							//fmt.Println(u)
							if err = ModifyADUser(l, u, usrOldDN); err == nil {
								sam := s.Username

								db.DeleteORM(orm, sam, u)

								for _, g := range e.GetAttributeValues("memberOf") {
									if !strings.Contains(g, "F_O365") {
										RemoveADGroupMember(l, baseDN, g, usrOldDN)
									}
								}
								//fmt.Println(usrOldDN)
								usrNewDn := strings.Replace(usrOldDN, "OU=Standard Users", "OU=Resigned Users", 1)
								//fmt.Println(usrNewDn)
								MoveADObject(l, usrOldDN, usrNewDn)

								if common.ConvertDatetimeToDate(s.StartDate) < "2020-10-01" && strings.Contains(s.Username, ".") {
									subj := fmt.Sprintf("Employee resigned. Please delete the canadiansolar mailbox - %s", s.Username)
									body := fmt.Sprintf("There is an employee resigned whose start date is before 2020-10-01:\n\nName: %s\nUser ID: %s\nSite: %s\nJob level: %s\nStart date: %s\nLast work date: %s", s.Username, s.PersonIDExternal, s.Site, s.JobLevel.String, common.ConvertDatetimeToDate(s.StartDate), common.ConvertDatetimeToDate(s.LastDateWorked))
									output.SendMail(output.MailReport, subj, body, "plain", "", output.MailAdmin, output.MailHelp)
									//data := fmt.Sprintf(mapZohoData[env]["disable"], s.Username, s.Username, s.PersonIDExternal, s.Site, s.JobLevel.String, common.ConvertDatetimeToDate(s.StartDate), common.ConvertDatetimeToDate(s.LastDateWorked))
									//output.CreateTicket(mapZohoURL[env], mapZohoKey[env], data)
								}
							}
						}
					}
				}
			}
		}
	}
}

func SyncADUsers(hour int, sqlquery, typ, env, update string) {
	ldapURL := MapLDAP[env]["url"]
	ldapADM := MapLDAP[env]["adm"] + "," + MapLDAP[env]["baseDN"]
	ldapPWD := MapLDAP[env]["pwd"]

	if typ == "update" {
		if hour == 0 { // 全量
			sqlquery = queryUpdateAll
		} else {
			sqlquery = fmt.Sprintf(sqlquery, strconv.Itoa(hour), strconv.Itoa(hour), strconv.Itoa(hour))
		}
		sqlquery += siteFilter + personFilter
		//sqlquery = `select * from "UNI_ID".."View_SFAttrForAD" where site = 'YCPT'`
		if update != "" {
			sqlquery = fmt.Sprintf(`SELECT * from UNI_ID..View_SFAttrForAD where username = '%s'`, update)
		}
	} else if typ == "disable" {
		if hour == 0 {
			sqlquery = queryDisableAll
		} else {
			sqlquery = fmt.Sprintf(sqlquery, strconv.Itoa(hour))
		}
		sqlquery += siteFilter + personFilter
		//sqlquery = fmt.Sprintf(`SELECT * from UNI_ID..View_SFAttrForAD WHERE username = 'billy.zhou'`)
	}
	output.GenerateLog(nil, "SQL query", sqlquery, true)

	r := db.QueryMsSql(sqlquery)
	defer r.Close()

	count := db.QueryMsSqlRow(strings.Replace(sqlquery, "*", "COUNT(*)", 1))

	orm := db.ConnectORM()

	l := BindAD(ldapURL, ldapADM, ldapPWD)
	defer l.Close()

	b := pgb.NewBar("AD user "+typ, count)
	for r.Next() {
		b.Add()
		var s = new(db.SFView)
		r.Scan(&s.PersonIDExternal, &s.StartDate, &s.CreatedOn, &s.LastDateWorked, &s.Site, &s.CompanyEN, &s.CompanyCN, &s.AddressEN, &s.AddressCN, &s.EmplStatus, &s.ManagerID, &s.FirstName, &s.LastName, &s.Username, &s.EmailAddress, &s.PhoneNumber, &s.SignLang, &s.JobLevel, &s.Department, &s.DepartmentCN, &s.TitleEN, &s.TitleCN, &s.JobEffectDate, &s.PosEffectDate, &s.LastModDate)
		//fmt.Println(s.Username)
		SyncADUser(l, s, typ, env, orm)
	}
	fmt.Println()
}

func InformADUserGroupManagers(l *ldap.Conn, baseDN, uSAM, old, new, whatChanged string) {
	mgrMap := GetADUserGroupManager(l, baseDN, uSAM)
	for k, v := range mgrMap {
		if whatChanged == "site" {
			output.SendMail(output.MailReport, "User site changed - "+uSAM, "Dear "+k+" group manager:\n\n  Please be noticed that the site of the user ("+
				uSAM+") has been changed from "+old+" to "+new+". AD Group members may need modification.\n\n  Please feel free to contact help@csisolar.com if you have any problem. Thank you.", "plain", "", output.MailAdmin, v)
		} else if whatChanged == "dept" {
			output.SendMail(output.MailReport, "User department changed - "+uSAM, "Dear "+k+" group manager:\n\n  Please be noticed that the department of the user ("+
				uSAM+") has been changed from "+old+" to "+new+". AD Group members may need modification.\n\n  Please feel free to contact help@csisolar.com if you have any problem. Thank you.", "plain", "", output.MailAdmin, v)
		}
	}
}

func SendADOperationLog() {
	sqlCreate := `SELECT sam_account_name, employee_id, physical_delivery_office_name, title, department, user_principal_name, created_at, deleted_at FROM UNI_ID..ad_user_attrs WHERE created_at >= DATEADD(DAY, -7, GETDATE());`
	sqlDelete := `SELECT sam_account_name, employee_id, physical_delivery_office_name, title, department, user_principal_name, created_at, deleted_at FROM UNI_ID..ad_user_attrs WHERE deleted_at >= DATEADD(DAY, -7, GETDATE());`

	rowsCreate := db.QueryMsSql(sqlCreate)
	rowsDelete := db.QueryMsSql(sqlDelete)

	f := excelize.NewFile()
	index := f.NewSheet("created")
	index = f.NewSheet("deleted")
	f.DeleteSheet("Sheet1")

	f.SetCellValue("created", "A1", "username")
	f.SetCellValue("created", "B1", "employee_id")
	f.SetCellValue("created", "C1", "site")
	f.SetCellValue("created", "D1", "title")
	f.SetCellValue("created", "E1", "department")
	f.SetCellValue("created", "F1", "upn")
	f.SetCellValue("created", "G1", "created_at")

	f.SetCellValue("deleted", "A1", "username")
	f.SetCellValue("deleted", "B1", "employee_id")
	f.SetCellValue("deleted", "C1", "site")
	f.SetCellValue("deleted", "D1", "title")
	f.SetCellValue("deleted", "E1", "department")
	f.SetCellValue("deleted", "F1", "upn")
	f.SetCellValue("deleted", "G1", "deleted_at")

	rowCount := 2
	for rowsCreate.Next() {
		a := new(db.ADUserAttr)

		rowsCreate.Scan(&a.SAMAccountName, &a.EmployeeID, &a.PhysicalDeliveryOfficeName, &a.Title, &a.Department, &a.UserPrincipalName, &a.CreatedAt, &a.DeletedAt)

		f.SetCellValue("created", "A"+strconv.Itoa(rowCount), a.SAMAccountName)
		f.SetCellValue("created", "B"+strconv.Itoa(rowCount), a.EmployeeID)
		f.SetCellValue("created", "C"+strconv.Itoa(rowCount), a.PhysicalDeliveryOfficeName)
		f.SetCellValue("created", "D"+strconv.Itoa(rowCount), a.Title)
		f.SetCellValue("created", "E"+strconv.Itoa(rowCount), a.Department)
		f.SetCellValue("created", "F"+strconv.Itoa(rowCount), a.UserPrincipalName)
		f.SetCellValue("created", "G"+strconv.Itoa(rowCount), strings.Split(a.CreatedAt.String(), " ")[0])

		rowCount++
	}

	rowCount = 2
	for rowsDelete.Next() {
		a := new(db.ADUserAttr)

		rowsDelete.Scan(&a.SAMAccountName, &a.EmployeeID, &a.PhysicalDeliveryOfficeName, &a.Title, &a.Department, &a.UserPrincipalName, &a.CreatedAt, &a.DeletedAt)

		f.SetCellValue("deleted", "A"+strconv.Itoa(rowCount), a.SAMAccountName)
		f.SetCellValue("deleted", "B"+strconv.Itoa(rowCount), a.EmployeeID)
		f.SetCellValue("deleted", "C"+strconv.Itoa(rowCount), a.PhysicalDeliveryOfficeName)
		f.SetCellValue("deleted", "D"+strconv.Itoa(rowCount), a.Title)
		f.SetCellValue("deleted", "E"+strconv.Itoa(rowCount), a.Department)
		f.SetCellValue("deleted", "F"+strconv.Itoa(rowCount), a.UserPrincipalName)
		f.SetCellValue("deleted", "G"+strconv.Itoa(rowCount), strings.Split(a.DeletedAt.Time.String(), " ")[0])

		rowCount++
	}

	file := "adoplog" + time.Now().Format("20060102") + ".xlsx"
	f.SetActiveSheet(index)
	err := f.SaveAs(file)
	output.GenerateLog(err, "Generate Excel file", file, false)

	output.SendMail(output.MailReport, "AD user operation log "+time.Now().Format("20060102"), "", "plain", file, "", output.MailAdmin)
}
