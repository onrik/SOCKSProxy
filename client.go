package proxy

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

type Client interface {
	Dial(network, address string) (net.Conn, error)
}

func NewClient(proxy *url.URL, conf *SOCKSConf) (client Client, err error) {
	switch strings.ToUpper(proxy.Scheme) {
	case "SOCKS", "SOCKS4", "SOCKS4A":
		client = &socks4Client{proxy, conf}
	case "SOCKS5", "SOCKS5+TLS":
		client = &socks5Client{proxy, conf, conf.TLSConfig}
	default:
		err = fmt.Errorf("%s not supported", proxy.Scheme)
	}
	return
}

func splitHostPort(addr string) (host, port []byte, err error) {
	hostName, hostPort, err := net.SplitHostPort(addr)
	if err != nil {
		return
	}
	_port, err := strconv.ParseUint(hostPort, 10, 16)
	_portBuffer := make([]byte, 2)
	binary.BigEndian.PutUint16(_portBuffer, uint16(_port))
	return []byte(hostName), _portBuffer, err
}
