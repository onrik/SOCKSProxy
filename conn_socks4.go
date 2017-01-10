package proxy

import (
	"io"
)

type socks4Conn struct {
	*socksConn
}

func (conn *socks4Conn) Serve() (err error) {
	request, err := readSocks4Request(conn.localConn)
	if err != nil {
		conn.sendReply(request, socks4StatusRejected)
		return err
	}
	switch request.command {
	case commandConnect:
		err = conn.handleConnect(request)
	default:
		err = errCommandNotSupported
	}
	if err != nil {
		conn.sendReply(request, socks4StatusRejected)
	}
	return
}

func (conn *socks4Conn) handleConnect(request *socks4Request) (err error) {
	conn.sendReply(request, socks4StatusGranted)
	remoteConn, err := conn.Dial("tcp", request.Address())
	if err != nil {
		return err
	}
	go io.Copy(conn.localConn, remoteConn)
	go io.Copy(remoteConn, conn.localConn)
	return
}

func (conn *socks4Conn) sendReply(request *socks4Request, status byte) {
	response := &socks4Response{
		status:status,
		port:  make([]byte, 2),
		ip:    make([]byte, 4),
	}
	if request.IsSOCKS4A() {
		response.port = request.port
		response.ip = request.ip
	}
	conn.localConn.Write(response.ToPacket())
	if status != socks4StatusGranted {
		conn.localConn.Close()
	}
}
