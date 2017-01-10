package proxy

import (
	"net"
	"net/url"
	"strings"
)

type socks4Client struct {
	proxy        *url.URL
	upstreamDial func(network, address string) (net.Conn, error)
}

func (c *socks4Client) Dial(_, address string) (remoteConn net.Conn, err error) {
	remoteConn, err = c.upstreamDial("tcp", c.proxy.Host)
	if err != nil {
		return
	}
	host, port, err := splitHostPort(address)
	if err != nil {
		return
	}
	request := &socks4Request{
		command:commandConnect,
		port:   port,
		ip:     []byte{0, 0, 0, 1},
		userId: []byte{},
		fqdn:   host,
	}
	if c.isSOCKS4() {
		request.ip, err = lookupIP(string(host))
		if err != nil {
			return
		}
	}
	response, err := sendReceive(remoteConn, request.ToPacket())
	if err != nil {
		return
	}
	if len(response) != 8 {
		err = errors.New("Server does not respond properly.")
		return
	}
	switch response[1] {
	case socks4StatusRejected:
		err = errors.New("Socks connection request rejected or failed.")
	case 92:
		err = errors.New("Socks connection request rejected becasue SOCKS server cannot connect to identd on the client.")
	case 93:
		err = errors.New("Socks connection request rejected because the client program and identd report different user-ids.")
	default:
		err = errors.New("Socks connection request failed, unknown error.")
	}
	return
}

func (c *socks4Client) isSOCKS4() bool {
	return strings.ToLower(c.proxy.Scheme) == "socks4"
}
