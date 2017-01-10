package proxy

import "errors"

const (
	socks4 byte = iota + 4
	socks5

	commandConnect      byte = iota + 1
	commandPortBinding
	commandUDPAssociate

	socks4StatusGranted  byte = 90
	socks4StatusRejected

	socks5AddressTypeIPv4 byte = iota + 1
	_
	socks5AddressTypeFQDN
	socks5AddressTypeIPv6

	socks5StatusSucceeded               byte = iota
	socks5StatusGeneral
	socks5StatusNotAllowed
	socks5StatusNetworkUnreachable
	socks5StatusHostUnreachable
	socks5StatusConnectionRefused
	socks5StatusTTLExpired
	socks5StatusCommandNotSupported
	socks5StatusAddressTypeNotSupported

	socks5AuthMethodNo       byte = iota
	_
	socks5AuthMethodPassword
)

var (
	errVersionError           = errors.New("version error")
	errUnsupportedCommand     = errors.New("unsupported command")
	errUnsupportedAuthMethod  = errors.New("unsupported authentication method")
	errUnsupportedAddressType = errors.New("unsupported address type")
)
