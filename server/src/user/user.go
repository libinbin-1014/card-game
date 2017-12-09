/* ====================================================================*/
/* Copyright (c) 2017.  All rights reserved.                           */
/* Author:     libinbin_1014@sina.com                                  */
/* Date :      2017/12/09                                              */
/* ====================================================================*/
package user

import (
	"errors"
	"sync"
)

type User struct {
	Num  int
	Name string
	Age  int
	Exp  int
	Pwd  string
	rw   sync.RWMutex
}

var UserMap map[int]User

func GetAccountInfo(num int) (User, error) {

	//get user info from map
	if v, ok := UserMap[num]; ok {
		return v, nil
	} else {
		return User{}, errors.New("Not Find the user")
	}
}

func ModifyAccount() {

}

func DeleteAccount() {

}

func CreateAccount() {

}
