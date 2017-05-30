package db

import (
	"time"

	"git.nulana.com/bobrnor/sqlsugar.git"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	if err := sqlsugar.Open("mysql", "root@unix(/tmp/mysql.sock)/battleship?parseTime=true"); err != nil {
		zap.S().Fatalw("can't open mysql connection",
			"err", err,
		)
		return
	}

	sqlsugar.SetMaxOpenConns(1)
	sqlsugar.SetMaxIdleConns(0)
	sqlsugar.SetConnMaxLifetime(10 * time.Second)
}
