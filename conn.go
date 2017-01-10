package proxy

import (
	"bufio"
	"io"
	"net"
)

type SOCKSConf struct {
	Auth        func(username, password string) bool
	Dial        func(network, address string) (net.Conn, error)
	HandleError func(error)
}

func Serve(listener net.Listener, conf *SOCKSConf) {
	if conf.HandleError == nil {
		conf.HandleError = func(_ error) {}
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			conf.HandleError(err)
			continue
		}
		version, err := bufio.NewReader(conn).ReadByte()
		if err != nil {
			return
		}
		switch version {
		case socks4version:
			if conf.Auth != nil {
				return
			}
			socksConn := socks4Conn{conn, conf}
			err = socksConn.Serve()
		case socks5version:
			socksConn := socks5Conn{conn, conf}
			err = socksConn.Serve()
		}
		if err != nil {
			conf.HandleError(err)
			continue
		}
	}
}

func FastMatch(r io.Reader) bool {
	header := make([]byte, 1, 1)
	if _, err := r.Read(header); err != nil {
		return false
	}
	return header[0] == 4 || header[0] == 5
}
