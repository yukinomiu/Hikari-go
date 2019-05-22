package hikariclient

import "errors"

const (
	// version
	socks5Ver byte = 5

	// auth method
	socks5MethodNoAuth        byte = 0
	socks5NoAcceptableMethods byte = 0xFF

	// command
	socks5CommandConnect byte = 1
	// socks5CommandBind         byte = 2
	// socks5CommandUdpAssociate byte = 3

	// rsv
	socks5Rsv byte = 0

	// address type
	socks5AddressTypeIpv4       byte = 1
	socks5AddressTypeDomainName byte = 3
	socks5AddressTypeIpv6       byte = 4

	// reply
	socks5ReplyOk                   byte = 0
	socks5ReplyGeneralServerFailure byte = 1
	socks5ReplyConnectionNotAllowed byte = 2
	socks5ReplyNetworkUnreachable   byte = 3
	socks5ReplyHostUnreachable      byte = 4
	// socks5ReplyConnectionRefused       byte = 5
	// socks5ReplyTtlExpired              byte = 6
	socks5ReplyCommandNotSupported byte = 7
	// socks5ReplyAddressTypeNotSupported byte = 8
)

var (
	// socks errors
	socksVersionNotSupportedErr  = errors.New("socks version not supported")
	socksMethodsNotAcceptableErr = errors.New("socks method not acceptable")
	socksCmdNotSupportedErr      = errors.New("socks5 command not supported")
	socksAdsTypeNotSupportedErr  = errors.New("socks5 address type not supported")

	// others
	badSocksAuthReqErr = errors.New("bad socks5 auth request")
	badSocksReqErr     = errors.New("bad socks5 request")
)
