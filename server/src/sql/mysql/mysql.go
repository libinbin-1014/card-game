/* ====================================================================*/
/* Copyright (c) 2017.  All rights reserved.                           */
/* Author:     libinbin_1014@sina.com                                  */
/* Date :      2017/12/09                                              */
/* ====================================================================*/
package mysql

import (
	"database/sql"
	"fmt"
	"sync"

	"../../conf"
	"../../log"
	"../../user"
	_ "github.com/go-sql-driver/mysql"
)

var Tlog = golog.GetLogHaddle()

type db struct {
	db   *sql.DB
	lock sync.Mutex
}

var DbHd = db{}

func MysqlInit() error {
	// 数据源字符串：用户名:密码@协议(地址:端口)/数据库?参数=参数值
	sqlStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", conf.DbUser, conf.DbPwd, conf.DbHost, conf.Dbname)
	t, err := sql.Open("mysql", sqlStr)
	if err != nil {
		return err
	}
	DbHd.db = t
	return nil
}

func MysqlUninit() {
	DbHd.db.Close()
	DbHd.lock.Unlock()
}

func AddAccount(u user.User) error {
	sql := fmt.Sprintf("insert into account (Num, Name, Pwd, Age, Exp) values(%d, \"%s\", \"%s\", %d, %d);",
		u.Num, u.Name, u.Pwd, u.Age, u.Exp)
	_, err := DbHd.db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func DelAccount(num int) error {
	sql := fmt.Sprintf("delete from account where Num = %d;", num)
	_, err := DbHd.db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func ModifyAccount(num int, u user.User) error {
	sql := fmt.Sprintf("update account set Name = '%s', Pwd = '%s', Age = %d, Exp = %d where Num = %d;",
		u.Name, u.Pwd, u.Age, u.Exp, u.Num)
	_, err := DbHd.db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func GetAccount(num int) (user.User, error) {
	sql := fmt.Sprintf("select  Name, Pwd, Age, Exp from account where Num = %d;", num)
	rows := DbHd.db.QueryRow(sql)

	var oneuser user.User

	rows.Scan(&oneuser.Name, &oneuser.Pwd, &oneuser.Age, &oneuser.Exp)

	oneuser.Num = num
	return oneuser, nil
}
