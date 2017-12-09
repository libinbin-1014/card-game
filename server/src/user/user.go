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
