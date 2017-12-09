/* ====================================================================*/
/* Copyright (c) 2017.  All rights reserved.                           */
/* Author:     libinbin_1014@sina.com                                  */
/* Date :      2017/12/09                                              */
/* ====================================================================*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"

	"./api/FileApi"
	"./conf"
	"./daemon"
	"./log"
	"./socket"
	"./sql/mysql"
	//	"./user"
)

var (
	IsDaemon bool
	IsHelp   bool
	IsStop   bool
	GetVer   bool
	Tlog     *golog.Logger
)

func ConfInit() {
	conf.GetConfInit()
}

func SocketInit() {
	defer Tlog.Infoln("Socket init success")
	go socket.Init()
}

func LogInit() {
	golog.LogInit()
	Tlog = golog.GetLogHaddle()
	Tlog.Infoln("Log init success")
}

func KillServer() {
	if FileApi.ChkExist(conf.PidFile) {
		os.Remove(conf.PidFile)
	}
}

func ArgsInit() {
	flag.BoolVar(&IsDaemon, "d", false, "input d to tarans daemon")
	flag.BoolVar(&IsHelp, "h", false, "print help info")
	flag.BoolVar(&IsStop, "s", false, "kill the server")
	flag.BoolVar(&GetVer, "v", false, "print the server version")

	flag.Parse()
	if IsHelp {
		flag.Usage()
	}
	if GetVer {
		fmt.Printf("the version is %s\n", conf.Version)
	}
	if IsStop {
		KillServer()
	}
	if IsHelp || IsStop || GetVer {
		os.Exit(0)
	}
}

func daemonInit() {
	ret := daemon.Daemon(0, 1)
	if ret == 0 {
		Tlog.Infoln("daemon process success")
	}
}

func ProcessInit() {
	pid := os.Getpid()

	// chk pid file
	if FileApi.ChkExist(conf.PidFile) != true {
		fp, err := os.Create(conf.PidFile)
		if err != nil {
			fmt.Println("create PidFile error ", err)
			os.Exit(1)
		}
		fp.Write([]byte(strconv.Itoa(pid)))
		fp.Close()
		return
	}
	Tlog.Warnln("the server is runing already, exit")
	os.Exit(1)
}

func EnvInit() {
	if conf.LogPath == "" || conf.LogFileName == "" {
		fmt.Println("get conf error")
		os.Exit(1)
	}

	//chk log dir
	//check the dir exist
	_, err := os.Stat(conf.LogPath)
	if FileApi.ChkExist(conf.LogPath) != true {
		if err := os.MkdirAll(conf.LogPath, 0777); err != nil {
			fmt.Println("MkdirAll error")
			os.Exit(1)
		}
	}

	//check the file exist
	if FileApi.ChkExist(conf.LogFileName) != true {
		if _, err = os.Create(conf.LogFileName); err != nil {
			fmt.Println("create logfile error ", err)
			os.Exit(1)
		}
	}

}

func SqlInit() {
	err := mysql.MysqlInit()
	if err != nil {
		Tlog.Errorln("mysql init err:", err)
		return
	}
	Tlog.Infoln("mysql init ok")
}

func main() {
	defer KillServer()
	ConfInit()
	EnvInit()
	ArgsInit()
	LogInit()
	//ProcessInit()
	SocketInit()
	SqlInit()

	var wg sync.WaitGroup
	wg.Add(1)
	if IsDaemon {
		daemonInit()
	}
	wg.Wait()
}
