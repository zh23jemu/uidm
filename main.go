package main

import (
	"flag"
	"time"
	"uidm/ad"
	"uidm/output"
	"uidm/sf"
)

const (
	sqlqueryNew     = `SELECT * from UNI_ID..View_SFAttrForAD WHERE username != 'null' AND lastDateWorked < startDate AND (lastModDate > DATEADD(hour,-%s,GETDATE()) OR jobEffectDate > DATEADD(hour,-%s,GETDATE()) OR posEffectDate > DATEADD(hour,-%s,GETDATE()))`
	sqlqueryDisable = `SELECT * FROM UNI_ID..View_SFAttrForAD WHERE username != 'null' AND lastDateWorked >= startDate AND lastDateWorked <= GETDATE() AND lastModDate >  DATEADD(hour,-%s,GETDATE())`
)

var (
	opLogFlag    = 0
	isRemindSent = false

	// Augments
	intervalHour = flag.Int("hour", 1, "How long of past data will be synced in hour? (int) Default is 1 hour.")
	env          = flag.String("env", "dev", "Which environment to be run? (dev/prd) Default is dev.")
	update       = flag.String("update", "", "To update specific username.")
	//isUpdate     = flag.String("update", "n", "Whether AD update operation will be implemented. (y/n) Default is n.")
)

func main() {
	flag.Parse()
	//var dayBeforeRun, dayAfterRun int

	//for {
	output.PrepareLog()
	dayBeforeRun := time.Now().Day()

	if dayBeforeRun == 5 && *intervalHour == 0 {
		//if !isRemindSent {
		output.SendMonthlyRemind()
	}
	//isRemindSent = true
	//}
	//} else {
	//	isRemindSent = false
	//}

	//if dayAfterRun != 0 && dayBeforeRun != dayAfterRun { //等于0代表是第一次运行，两个时间不相等代表是新的一天
	//	sf.SyncPosition(0)
	//	sf.SyncPerPersonal(0)
	//	sf.SyncEmpEmployment(0)
	//	sf.SyncFOCompany(0)
	//	sf.SyncEmpJob(0)
	//	sf.SyncFODepartment(0)
	//	sf.SyncPerEmail(0)
	//	sf.SyncPerNationalID(0)
	//	sf.SyncPerPhone(0)
	//	sf.SyncUser(0)
	//}

	sf.SyncPosition(*intervalHour)
	sf.SyncPerPersonal(*intervalHour)
	sf.SyncEmpEmployment(*intervalHour)
	sf.SyncFOCompany(*intervalHour)
	sf.SyncEmpJob(*intervalHour)
	sf.SyncFODepartment(*intervalHour)
	sf.SyncPerEmail(*intervalHour)
	sf.SyncPerNationalID(*intervalHour)
	sf.SyncPerPhone(*intervalHour)
	sf.SyncUser(*intervalHour)
	sf.SyncServiceCompany()

	ad.SyncADUsers(*intervalHour, sqlqueryNew, "update", *env, *update)
	ad.SyncADUsers(*intervalHour, sqlqueryDisable, "disable", *env, *update)
	if *intervalHour != 0 {
		output.SendLogMail()
	}

	//dayAfterRun = time.Now().Day()

	//interval, _ := time.ParseDuration("30m")
	//NextRun := time.Now().Add(interval).Format("2006-01-02 15:04")
	//fmt.Println("Next sync will be started at: " + NextRun)
	//
	//time.Sleep(30 * time.Minute)
	//}
}
