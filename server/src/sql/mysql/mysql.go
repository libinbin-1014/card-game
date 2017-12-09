package mysql

import (
	"database/sql"
	"fmt"
	"sync"

	"../../conf"
	"../../user"
	_ "github.com/go-sql-driver/mysql"
)

type db struct {
	db   sql.DB
	lock sync.Mutex
}

var DbHd = db{}

func MysqlInit() {
	sqlStr := conf.DbUser + ":" + conf.DbPwd + "@(" + conf.DbHost + ")/" + conf.Dbname
	t, err := sql.Open("mysql", sqlStr)
	DbHd.db = *t
	fmt.Println(DbHd.db, err)
}

func MysqlUninit() {
	DbHd.db.Close()
	DbHd.lock.Unlock()
}

func Add(u user.User) error {
	sql := fmt.Sprintf("insert into account (Num, Name, Pwd, Age, Exp) values(%d, \"%s\", \"%s\", %d, %d);",
		u.Num, u.Name, u.Pwd, u.Age, u.Exp)
	fmt.Println(sql)
	return nil
}
