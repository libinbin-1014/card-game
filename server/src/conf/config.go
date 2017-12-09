/* ====================================================================*/
/* Copyright (c) 2017.  All rights reserved.                           */
/* Author:     libinbin_1014@sina.com                                  */
/* Date :      2017/12/09                                              */
/* ====================================================================*/
package conf

import (
	"bufio"
	"io"
	"os"
	"strings"
)

var (
	LogPath     string = "/var/log/code/"
	Version     string = "1.1"
	LogFileName string = "sys.log"
	PidFile     string = "/var/run/demo.pid"
	Port        string = "8888"
	DbHost      string
	DbUser      string
	DbPwd       string
	Dbname      string
)

const middle = ":"

type Config struct {
	Mymap  map[string]string
	strcet string
}

func (c *Config) initConfig(path string) {
	c.Mymap = make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		s := strings.TrimSpace(string(b))
		//fmt.Println(s)
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			c.strcet = strings.TrimSpace(s[n1+1 : n2])
			continue
		}

		if len(c.strcet) == 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}

		key := c.strcet + middle + frist
		c.Mymap[key] = strings.TrimSpace(second)
	}
}

func (c Config) read(node, key, defaultValue string) string {
	key = node + middle + key
	v, found := c.Mymap[key]
	if !found {
		return defaultValue
	}
	return v
}
func GetConfInit() {
	defaultPath := "./conf/server.conf"
	myConfig := new(Config)
	myConfig.initConfig(defaultPath)

	//log
	LogPath = myConfig.read("log", "logpath", LogPath)
	LogFileName = myConfig.read("log", "logfilename", LogFileName)
	LogFileName = LogPath + LogFileName

	//server
	Version = myConfig.read("server", "version", Version)
	PidFile = myConfig.read("server", "pidfile", PidFile)

	//socket
	Port = myConfig.read("socket", "port", Port)

	//mysql
	DbHost = myConfig.read("mysql", "dbhostip", DbHost)
	DbUser = myConfig.read("mysql", "dbusername", DbUser)
	DbPwd = myConfig.read("mysql", "dbpwd", DbPwd)
	Dbname = myConfig.read("mysql", "dbname", Dbname)
}
