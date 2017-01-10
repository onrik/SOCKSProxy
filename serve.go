package proxy

import (
	"io"
	"net"
)

type SOCKSConf struct {
	Auth        func(username, password string) bool
	Dial        func(network, address string) (net.Conn, error)
	HandleError func(error)
}

func Listen(listener net.Listener, conf *SOCKSConf) {
	if conf.HandleError == nil {
		conf.HandleError = func(_ error) {}
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			conf.HandleError(err)
			continue
		}
		socksConn := socksConn{conn, conf}
		err = socksConn.Serve()
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
