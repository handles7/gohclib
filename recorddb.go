package gohclib

import (
	"fmt"
	"time"

	"github.com/handles7/gohclib/cron"

	"github.com/handles7/gohclib/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func RunRecord(siid int, releasesiid int) {
	datasource := config.GetStr("mysqlpath")
	prod := config.GetStr("PROFILE")
	var sid int

	if prod == "prod" {
		sid = releasesiid
	} else {
		sid = siid
	}

	if datasource != "" {
		c := cron.New()
		c.AddFunc("@every 30s", func() {
			RecordCenterCron(sid, datasource)

		})
		c.Start()
	}

}

func RecordCenterCron(siid int, datasource string) {

	dbc, err := sqlx.Open("mysql", datasource)

	if err != nil {
		fmt.Println(err.Error())
		return
	} else {
		defer dbc.Close()
	}

	dbc.Unsafe()
	dbc.SetMaxOpenConns(2)
	dbc.SetMaxIdleConns(1)

	sqlu := `update server_status set ssupdatetime=:ssupdatetime where siid=:siid`
	dbc.NamedExec(sqlu, map[string]interface{}{
		"ssupdatetime": time.Now().Format("2006-01-02 15:04:05"),
		"siid":         siid,
	})

}
