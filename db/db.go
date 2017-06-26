package db

import (
	"time"

	"git.nulana.com/bobrnor/sqlsugar.git"

	"log"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	if err := sqlsugar.Open("mysql", "root@tcp(172.25.0.2:3306)/battleship?parseTime=true"); err != nil {
		log.Fatalf("can't open mysql connection",
			"err", err,
		)
		return
	}

	sqlsugar.SetMaxOpenConns(1)
	sqlsugar.SetMaxIdleConns(0)
	sqlsugar.SetConnMaxLifetime(10 * time.Second)
}
