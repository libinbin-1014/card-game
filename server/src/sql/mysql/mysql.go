/* ====================================================================*/
/* Copyright (c) 2017.  All rights reserved.                           */
/* Author:     libinbin_1014@sina.com                                  */
/* Date :      2017/12/09                                              */
/* ====================================================================*/
package mysql

import (
	"database/sql"
	"errors"
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

	// get all account in user.UserMap
	// that can quick get account
	rows, err := t.Query("select Num,Name,Pwd,Age,Exp from account;")
	if err != nil {
		return err
	}

	user.MapRsync.Lock()
	defer user.MapRsync.Unlock()
	for rows.Next() {
		var oneUser user.User
		err = rows.Scan(&oneUser.Num, &oneUser.Name, &oneUser.Pwd, &oneUser.Age, &oneUser.Exp)
		if err != nil {
			return err
		}
		user.UserMap[oneUser.Num] = oneUser
	}
	return nil
}

func MysqlUninit() {
	DbHd.db.Close()
	DbHd.lock.Unlock()
}

func AddAccount(u user.User) error {
	user.MapRsync.Lock()
	defer user.MapRsync.Unlock()
	if _, ok := user.UserMap[u.Num]; ok {
		return ModifyAccount(u.Num, u)
	}

	sql := fmt.Sprintf("insert into account (Num, Name, Pwd, Age, Exp) values(%d, \"%s\", \"%s\", %d, %d);",
		u.Num, u.Name, u.Pwd, u.Age, u.Exp)
	_, err := DbHd.db.Exec(sql)
	if err != nil {
		return err
	}
	user.UserMap[u.Num] = u
	Tlog.Debugln("AddAccount success:", u)
	return nil
}

func DelAccount(num int) error {
	user.MapRsync.Lock()
	defer user.MapRsync.Unlock()
	if _, ok := user.UserMap[num]; ok {
		delete(user.UserMap, num)
	} else {
		return nil
	}

	sql := fmt.Sprintf("delete from account where Num = %d;", num)
	_, err := DbHd.db.Exec(sql)
	if err != nil {
		return err
	}
	Tlog.Debugln("DelAccount success: the account num is ", num)
	return nil
}

func ModifyAccount(num int, t_user user.User) error {
	user.MapRsync.Lock()
	defer user.MapRsync.Unlock()
	u := t_user
	var one user.User
	one, ok := user.UserMap[num]
	if !ok {
		return errors.New("modify failed, can not find the user")
	}

	//chk the changed value
	if u.Num == one.Num && u.Name == one.Name && u.Pwd == one.Pwd && u.Age == one.Age && u.Exp == one.Exp {
		return nil
	}
	if u.Num == 0 {
		u.Num = one.Num
	}
	if u.Name == "" {
		u.Name = one.Name
	}
	if u.Pwd == "" {
		u.Pwd = one.Pwd
	}
	if u.Age == 0 {
		u.Age = one.Age
	}
	if u.Exp == 0 {
		u.Exp = one.Exp
	}

	// update the sql
	sql := fmt.Sprintf("update account set Name = '%s', Pwd = '%s', Age = %d, Exp = %d where Num = %d;",
		u.Name, u.Pwd, u.Age, u.Exp, u.Num)
	_, err := DbHd.db.Exec(sql)
	if err != nil {
		return err
	}
	user.UserMap[num] = u
	Tlog.Debugln("ModifyAccount success, form ", one, " to ", u)
	return nil
}

func GetAccount(num int) (user.User, error) {
	user.MapRsync.RLock()
	defer user.MapRsync.RUnlock()
	var one user.User
	if one, ok := user.UserMap[num]; ok {
		return one, nil
	}
	return one, errors.New("not found the account")
}
