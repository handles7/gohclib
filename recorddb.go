package gohclib

import (
	"fmt"
	"gohclib/config"
	"gohclib/cron"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func RunRecord(siid int) {
	datasource := config.GetStr("mysqlpath")

	if datasource != "" {
		c := cron.New()
		c.AddFunc("@every 30s", func() { RecordCenterCron(siid, datasource) })
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
