package hikariserver

import (
	"hikari-go/hikaricommon"
	"net"
)

var queueCh = make(chan *context, ctxQueueSize)

type context struct {
	plainBuf     []byte
	encryptedBuf []byte
	clientConn   *net.Conn
	targetConn   *net.Conn
	crypto       *hikaricommon.Crypto
}

func (c *context) Close() {
	cc, tc := c.clientConn, c.targetConn

	if cc != nil {
		(*cc).Close()
	}

	if tc != nil {
		(*tc).Close()
	}
}

func (c *context) GetPlainBuf() []byte {
	return c.plainBuf
}

func (c *context) GetEncryptedBuf() []byte {
	return c.encryptedBuf
}

func (c *context) GetPlainConn() *net.Conn {
	return c.targetConn
}

func (c *context) GetEncryptedConn() *net.Conn {
	return c.clientConn
}

func (c *context) GetCrypto() *hikaricommon.Crypto {
	return c.crypto
}

func getContext() *context {
	var ctx *context

	select {
	case ctx = <-queueCh:
	default:
		ctx = &context{plainBuf: make([]byte, ctxBufSize), encryptedBuf: make([]byte, ctxBufSize)}
	}

	return ctx
}

func returnContext(ctx *context) {
	ctx.clientConn = nil
	ctx.targetConn = nil
	ctx.crypto = nil

	select {
	case queueCh <- ctx:
	default:
	}
}
