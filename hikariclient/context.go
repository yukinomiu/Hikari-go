package hikariclient

import (
	"hikari-go/hikaricommon"
	"net"
)

var queueCh = make(chan *context, ctxQueueSize)

type context struct {
	plainBuf     []byte
	encryptedBuf []byte
	localConn    *net.Conn
	serverConn   *net.Conn
	crypto       *hikaricommon.Crypto
}

func (c *context) Close() {
	lc, sc := c.localConn, c.serverConn

	if lc != nil {
		_ = (*lc).Close()
	}

	if sc != nil {
		_ = (*sc).Close()
	}
}

func (c *context) GetPlainBuf() []byte {
	return c.plainBuf
}

func (c *context) GetEncryptedBuf() []byte {
	return c.encryptedBuf
}

func (c *context) GetPlainConn() *net.Conn {
	return c.localConn
}

func (c *context) GetEncryptedConn() *net.Conn {
	return c.serverConn
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
	ctx.localConn = nil
	ctx.serverConn = nil
	ctx.crypto = nil

	select {
	case queueCh <- ctx:
	default:
	}
}
