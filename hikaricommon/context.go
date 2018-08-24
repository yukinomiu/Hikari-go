package hikaricommon

import "net"

type Context interface {
	Close()
	GetPlainBuf() []byte
	GetEncryptedBuf() []byte
	GetPlainConn() *net.Conn
	GetEncryptedConn() *net.Conn
	GetCrypto() *Crypto
}
