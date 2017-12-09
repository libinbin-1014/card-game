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

func CreateAccount(args []string) error {
	Tlog.Infoln("register a new account begin,", args)

	if string(args[0]) != "register" {
		return errors.New("CreateAccount para check err")
	}

	num, _ := strconv.Atoi(args[1])
	name := string(args[2])
	age, _ := strconv.Atoi(args[3])
	pwd := string(args[4])
	exp := 1000

	// chk the use exist
	oneuser, err := mysql.GetAccount(num)
	fmt.Println(oneuser)
	if err != nil {
		Tlog.Errorln("get account info err, ", err, args)
		return err
	}
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
		fmt.Println(tokens[0])
		switch string(tokens[0]) {
		case "login":
			fmt.Println("client is login")
		case "register":
			CreateAccount(tokens)
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
