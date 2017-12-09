package socket

import (
	"io"
	"net"

	"../conf"
	"../log"
)

var Tlog = golog.GetLogHaddle()

const BufLength = 218

func handleTcpCli(conn net.Conn) {
	defer conn.Close()
	Tlog.Infoln("accept a connect from:", conn.RemoteAddr())

	conn.Write([]byte("welcome to this room"))
	for {
		data := make([]byte, 0)
		buf := make([]byte, BufLength)
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
		Tlog.Debugln("recv the client data:", string(data))
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
