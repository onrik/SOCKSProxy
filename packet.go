package proxy

import (
	"bufio"
	"encoding/binary"
	"net"
	"strconv"
)

type socks4Request struct {
	command byte
	port    []byte
	ip      []byte
	userId  []byte
	fqdn    []byte
}

func readSocks4Request(conn net.Conn) (request *socks4Request, err error) {
	reader := bufio.NewReader(conn)
	request = &socks4Request{}
	request.command, err = reader.ReadByte()
	if err != nil {
		return
	}
	request.port, err = reader.ReadBytes(2)
	if err != nil {
		return
	}
	request.ip, err = reader.ReadBytes(4)
	if err != nil {
		return
	}
	request.userId, err = reader.ReadSlice(0)
	if err != nil {
		return
	}
	if !request.IsSOCKS4A() {
		return
	}
	request.fqdn, err = reader.ReadSlice(0)
	if err != nil {
		return
	}
	return
}

func (request *socks4Request) IsSOCKS4A() bool {
	ip := request.ip
	return ip[0] == 0 && ip[1] == 0 && ip[2] == 0 && ip[3] != 0
}

func (request *socks4Request) Address() string {
	var host, port string
	if request.IsSOCKS4A() {
		host = string(request.fqdn)
	} else {
		host = net.IP(request.ip).String()
	}
	port = strconv.Itoa(int(binary.BigEndian.Uint16(request.ip)))
	return net.JoinHostPort(host, port)
}

func (request *socks4Request) ToPacket() []byte {
	packet := []byte{SOCKS4Version, request.command}
	packet = append(packet, request.port...)
	packet = append(packet, request.ip...)

	packet = append(append(packet, request.userId...), 0)
	if request.IsSOCKS4A() {
		packet = append(append(packet, request.fqdn...), 0)
	}

	return packet
}

type socks5Initial struct {
	version byte
	methods []byte
}

func (request *socks5Initial) ToPacket() []byte {
	packet := []byte{
		request.version,
		byte(len(request.methods)),
	}
	packet = append(packet, request.methods...)
	return packet
}

type socks5Request struct {
	version  byte
	command  byte
	addrType byte
	addr     []byte
	port     []byte
}

func (request *socks5Request) ToPacket() []byte {
	packet := []byte{
		request.version,
		request.command,
		0x00,
		request.addrType,
	}
	packet = append(packet, request.addr...)
	packet = append(packet, request.port...)
	return packet
}

func (request *socks5Request) Address() string {
	var host string
	switch request.addrType {
	case socks5AddressTypeIPv4, socks5AddressTypeIPv6:
		host = net.IP(request.addr).String()
	case socks5AddressTypeFQDN:
		host = string(request.addr)
	}
	port := strconv.Itoa(int(binary.BigEndian.Uint16(request.port)))
	return net.JoinHostPort(host, port)
}

func readSocks5Request(conn net.Conn) (request *socks5Request, err error) {
	reader := bufio.NewReader(conn)
	request = &socks5Request{}
	request.version, err = reader.ReadByte()
	if err != nil {
		return
	}
	request.command, err = reader.ReadByte()
	if err != nil {
		return
	}
	_, err = reader.ReadByte()
	if err != nil {
		return
	}
	request.addrType, err = reader.ReadByte()
	if err != nil {
		return
	}
	switch request.addrType {
	case socks5AddressTypeIPv4:
		request.addr, err = reader.ReadBytes(net.IPv4len)
		if err != nil {
			return
		}
	case socks5AddressTypeIPv6:
		request.addr, err = reader.ReadBytes(net.IPv6len)
		if err != nil {
			return
		}
	case socks5AddressTypeFQDN:
		length, err := reader.ReadByte()
		if err != nil {
			return
		}
		request.addr, err = reader.ReadBytes(length)
		if err != nil {
			return
		}
	default:
		conn.Write(buildSocks5Reply(request, socks5StatusAddressTypeNotSupported))
		conn.Close()
		return nil, errUnsupportedAddressType
	}
	request.port, err = reader.ReadBytes(2)
	return
}
