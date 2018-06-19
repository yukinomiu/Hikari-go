package hikariclient

import (
	"bytes"
	"hikari-go/hikaricommon"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

// struct
type authReq struct {
	methods []byte
}

type socksReq struct {
	cmd        byte
	adsType    byte
	adsAndPort []byte
}

func startSocksServer() {
	if cfg.SocksPort == 0 {
		log.Println("socks server disabled")
		return
	}

	// init listener
	socksAds := cfg.SocksAddress
	socksPort := strconv.Itoa(int(cfg.SocksPort))
	socksListenAds := net.JoinHostPort(socksAds, socksPort)

	listener, err := net.Listen("tcp", socksListenAds)
	if err != nil {
		log.Fatalf("socks server listen on address '%v' err: %v\n", socksListenAds, err)
	}
	defer listener.Close()

	// accept
	log.Printf("socks server listen on address '%v'\n", socksListenAds)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept socks request err: %v\n", err)
			continue
		}

		go handleSocksConnection(&conn)
	}
	log.Println("socks server stop")
}

func handleSocksConnection(conn *net.Conn) {
	// init context
	ctx := &context{}
	defer func() {
		ctx.Close()

		if err := recover(); err != nil {
			log.Printf("unexpected err: %v\n", err)
		}
	}()

	// socks handshake
	if err := socksHandshake(ctx, conn); err != nil {
		log.Printf("socks handshake err: %v\n", err)
		return
	}

	// switch
	hikaricommon.Switch(ctx.localConn, ctx.serverConn, ctx.crypto, switchBufSize, switchTimeout)
}

func socksHandshake(ctx *context, conn *net.Conn) error {
	ctx.localConn = conn

	// deadline
	deadline := time.Now().Add(time.Second * handshakeTimeout)

	// set local connection timeout
	if err := (*conn).SetDeadline(deadline); err != nil {
		return err
	}

	// read auth req
	var authReq *authReq
	if req, err := readAuthReq(ctx.localConn); err != nil {
		return err
	} else {
		authReq = req
	}

	// reply auth req
	if err := replyAuthReq(ctx.localConn, authReq); err != nil {
		return err
	}

	// read socks req
	var socksReq *socksReq
	if req, err := readSocksReq(ctx.localConn); err != nil {
		return err
	} else {
		socksReq = req
	}

	// init crypto
	if c, err := hikaricommon.NewAESCrypto(secretKey, nil); err != nil {
		return err
	} else {
		crypto := interface{}(c).(hikaricommon.Crypto)
		ctx.crypto = &crypto
	}

	// send hikari req
	var hikariAdsType byte
	if t, err := toHikariAdsType(socksReq.adsType); err != nil {
		return err
	} else {
		hikariAdsType = t
	}

	if c, err := sendHikariReq(ctx.crypto, hikariAdsType, socksReq.adsAndPort); err != nil {
		ctx.serverConn = c
		return err
	} else {
		ctx.serverConn = c

		// deadline
		if err := (*c).SetDeadline(deadline); err != nil {
			return err
		}
	}

	// read hikari rsp
	var hikariRsp *hikariRsp
	if rsp, err := readHikariRsp(ctx.crypto, ctx.serverConn); err != nil {
		return err
	} else {
		hikariRsp = rsp
	}

	// reply socks req
	if err := replySocksReq(ctx.localConn, hikariRsp); err != nil {
		return err
	}

	return nil
}

func readAuthReq(conn *net.Conn) (*authReq, error) {
	c := *conn

	// init buffer
	buf := make([]byte, socksAuthReqBufSize)

	// read
	n, err := io.ReadAtLeast(c, buf, 2)
	if err != nil {
		return nil, err
	}

	// ver
	ver := buf[0]
	if ver != socks5Ver {
		return nil, socksVersionNotSupportedErr
	}

	// methods
	expectLen := int(buf[1]) + 2

	if n == expectLen {
	} else if n < expectLen {
		if expectLen > socksAuthReqBufSize {
			return nil, badSocksAuthReqErr
		}

		b := buf[n:expectLen]
		if _, err := io.ReadFull(c, b); err != nil {
			return nil, err
		}
	} else if n > expectLen {
		return nil, badSocksAuthReqErr
	}

	return &authReq{buf[2:expectLen]}, nil
}

func replyAuthReq(conn *net.Conn, req *authReq) error {
	c := *conn

	if !bytes.Contains(req.methods, []byte{0}) {
		if _, err := c.Write([]byte{socks5Ver, socks5NoAcceptableMethods}); err != nil {
			return err
		}

		return socksMethodsNotAcceptableErr
	}

	if _, err := c.Write([]byte{socks5Ver, socks5MethodNoAuth}); err != nil {
		return err
	}

	return nil
}

func readSocksReq(conn *net.Conn) (*socksReq, error) {
	c := *conn

	// init buffer
	buf := make([]byte, socksReqBufSize)

	// read
	n, err := io.ReadAtLeast(c, buf, 5)
	if err != nil {
		return nil, err
	}

	// ver
	ver := buf[0]
	if ver != socks5Ver {
		return nil, socksVersionNotSupportedErr
	}

	// command
	cmd := buf[1]
	if cmd != socks5CommandConnect {
		if _, err := c.Write([]byte{socks5Ver, socks5ReplyCommandNotSupported}); err != nil {
			return nil, err
		}
		return nil, socksCmdNotSupportedErr
	}

	// ignore rsv

	// address type
	adsType := buf[3]
	var expectLen int

	switch adsType {
	case socks5AddressTypeDomainName:
		expectLen = 6 + 1 + int(buf[4])
	case socks5AddressTypeIpv4:
		expectLen = 6 + net.IPv4len
	case socks5AddressTypeIpv6:
		expectLen = 6 + net.IPv6len
	default:
		return nil, socksAdsTypeNotSupportedErr
	}

	if n == expectLen {
	} else if n < expectLen {
		if expectLen > socksReqBufSize {
			return nil, badSocksReqErr
		}

		b := buf[n:expectLen]
		if _, err := io.ReadFull(c, b); err != nil {
			return nil, err
		}
	} else if n > expectLen {
		return nil, badSocksReqErr
	}

	adsAndPort := buf[4:expectLen]

	return &socksReq{
			cmd,
			adsType,
			adsAndPort},
		nil
}

func replySocksReq(conn *net.Conn, rsp *hikariRsp) error {
	c := *conn

	switch rsp.reply {
	case hikaricommon.HikariReplyOk:
		var socksAdsType byte
		if t, err := toSocksAdsType(rsp.bindAdsType); err != nil {
			return err
		} else {
			socksAdsType = t
		}

		// response
		buf := make([]byte, 4+len(rsp.bindAdsAndPort))
		buf[0] = socks5Ver
		buf[1] = socks5ReplyOk
		buf[2] = socks5Rsv
		buf[3] = socksAdsType
		copy(buf[4:], rsp.bindAdsAndPort)

		if _, err := c.Write(buf); err != nil {
			return err
		}

		if rsp.extraData != nil {
			if _, err := c.Write(rsp.extraData); err != nil {
				return err
			}
		}

		return nil

	case hikaricommon.HikariReplyVersionNotSupported:
		if _, err := c.Write([]byte{socks5Ver, socks5ReplyGeneralServerFailure}); err != nil {
			return err
		}
		return hikaricommon.HikariVerNotSupportedErr

	case hikaricommon.HikariReplyAuthFail:
		if _, err := c.Write([]byte{socks5Ver, socks5ReplyConnectionNotAllowed}); err != nil {
			return err
		}
		return hikaricommon.HikariAuthFailErr

	case hikaricommon.HikariAdsTypeNotSupported:
		if _, err := c.Write([]byte{socks5Ver, socks5ReplyGeneralServerFailure}); err != nil {
			return err
		}
		return hikaricommon.HikariAdsTypeNotSupportedErr

	case hikaricommon.HikariReplyDnsLookupFail:
		if _, err := c.Write([]byte{socks5Ver, socks5ReplyHostUnreachable}); err != nil {
			return err
		}
		return hikaricommon.HikariDnsLookupFailErr

	case hikaricommon.HikariReplyConnectTargetFail:
		if _, err := c.Write([]byte{socks5Ver, socks5ReplyNetworkUnreachable}); err != nil {
			return err
		}
		return hikaricommon.HikariConnectToTargetFailErr

	default:
		if _, err := c.Write([]byte{socks5Ver, socks5ReplyGeneralServerFailure}); err != nil {
			return err
		}
		return hikaricommon.BadHikariRspErr
	}
}

func toHikariAdsType(socksAdsType byte) (byte, error) {
	switch socksAdsType {
	case socks5AddressTypeDomainName:
		return hikaricommon.HikariAddressTypeDomainName, nil
	case socks5AddressTypeIpv4:
		return hikaricommon.HikariAddressTypeIpv4, nil
	case socks5AddressTypeIpv6:
		return hikaricommon.HikariAddressTypeIpv6, nil
	default:
		return 0, socksAdsTypeNotSupportedErr
	}
}

func toSocksAdsType(hikariAdsType byte) (byte, error) {
	switch hikariAdsType {
	case hikaricommon.HikariAddressTypeDomainName:
		return socks5AddressTypeDomainName, nil
	case hikaricommon.HikariAddressTypeIpv4:
		return socks5AddressTypeIpv4, nil
	case hikaricommon.HikariAddressTypeIpv6:
		return socks5AddressTypeIpv6, nil
	default:
		return 0, hikaricommon.HikariAdsTypeNotSupportedErr
	}
}
