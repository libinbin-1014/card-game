/* ====================================================================*/
/* Copyright (c) 2017.  All rights reserved.                           */
/* Author:     libinbin_1014@sina.com                                  */
/* Date :      2017/12/09                                              */
/* ====================================================================*/
package socket

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"../conf"
	"../log"
	"../sql/mysql"
	"../user"
)

var Tlog = golog.GetLogHaddle()

const BufLength = 218

func LoginAccount(args []string) (user.User, error) {
	var oneuser user.User

	if len(args) < 3 || string(args[0]) != "login" {
		return oneuser, errors.New("LoginAccount para check err\n please input login num pwd")
	}

	num, _ := strconv.Atoi(args[1])
	pwd := string(args[2])

	// get user info
	oneuser, _ = mysql.GetAccount(num)
	if oneuser.Name == "" {
		Tlog.Errorln("the user is not exist ,cat not login,", args)
		return oneuser, errors.New("user is not exist")
	}
	if pwd != oneuser.Pwd {
		Tlog.Errorln("the use pwd is wrong")
		return oneuser, errors.New("user pwd is wrong")
	}
	Tlog.Infoln("Login a account success")
	return oneuser, nil
}

func CreateAccount(args []string) error {
	Tlog.Infoln("register a new account begin,", args)

	if len(args) < 5 || string(args[0]) != "register" {
		return errors.New("CreateAccount para check err\n please input register num name age pwd")
	}

	num, _ := strconv.Atoi(args[1])
	name := string(args[2])
	age, _ := strconv.Atoi(args[3])
	pwd := string(args[4])
	exp := 1000

	// chk the use exist
	oneuser, err := mysql.GetAccount(num)
	if oneuser.Name != "" {
		Tlog.Errorln("the user is exist ,cat not register,", args)
		return errors.New("user is exist")
	}

	newuser := user.User{Num: num, Name: name, Age: age, Pwd: pwd, Exp: exp}
	err = mysql.AddAccount(newuser)
	if err != nil {
		return err
	}
	Tlog.Infoln("register a new account success")
	return nil
}

func handleTcpCli(conn net.Conn) {
	defer conn.Close()
	defer Tlog.Infoln(conn.RemoteAddr(), " connect closed")

	IsLogin := false
	var oneuser user.User
	conn.Write([]byte("welcome to this room"))
	for {
		data := make([]byte, 512)
		//buf := make([]byte, BufLength)
		n, err := conn.Read(data)
		if err != nil {
			Tlog.Errorln("recv the data error :", err)
			break
		}
		if n > 0 {
			data[n] = 0
		}
		reciveStr := string(data[0:n])
		/*
			for {
				n, err := conn.Read(buf)
				if err != nil && err != io.EOF {
					Tlog.Errorln("recv the data error:", err)
					break
				}
				data = append(data, buf[:n]...)
				if n != BufLength {
					break
				}
			}
		*/
		Tlog.Debugln("recv the client data:", reciveStr)

		tokens := strings.Split(reciveStr, " ")
		switch string(tokens[0]) {
		case "login":
			if !IsLogin {
				oneuser, err = LoginAccount(tokens)
				if err != nil {
					conn.Write([]byte(fmt.Sprintf("Login Failed: %s\n", err.Error())))
				} else {
					IsLogin = true
					conn.Write([]byte(fmt.Sprintf("Login Success\n user exp is %d\n", oneuser.Exp)))
				}
			} else {
				conn.Write([]byte(fmt.Sprintf("the user have been login\n the num is %d\n", oneuser.Num)))
			}
		case "register":
			if !IsLogin {
				err = CreateAccount(tokens)
				if err != nil {
					conn.Write([]byte(fmt.Sprintf("register failed:%s", err.Error())))
				} else {
					conn.Write([]byte("congratulation, register success"))
				}
			} else {
				conn.Write([]byte("have been login, do not need register"))
			}
		case "listuser":
			userNumList := []string{}
			user.MapRsync.RLock()
			for _, v := range user.UserMap {
				userNumList = append(userNumList, v.Name)
			}
			user.MapRsync.RUnlock()
			conn.Write([]byte(strings.Join(userNumList, ",")))
		default:
			fmt.Println("unknow the cmd")
		}

	}
}

func Init() {
	l, err := net.Listen("tcp", ":"+conf.Port)
	if err != nil {
		Tlog.Errorf("listen error: %s", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			Tlog.Errorf("listen error: %s", err)
			break
		}
		// start a new goroutine to handle
		// the new connection.
		go handleTcpCli(c)
	}
}
