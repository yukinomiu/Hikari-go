package hikariclient

import (
	"bufio"
	"encoding/binary"
	"hikari-go/hikaricommon"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	httpProxyOK   = []byte("HTTP/1.1 200 Connection Established\r\n\r\n")
	httpProxyFail = []byte("HTTP/1.1 502 Bad Gateway\r\n\r\n")
)

// struct
type httpConnectReq struct {
	adsAndPort []byte
}

func startHttpServer() {
	if cfg.HttpPort == 0 {
		log.Println("http server disabled")
		return
	}

	// init listener
	httpAds := cfg.HttpAddress
	httpPort := strconv.Itoa(int(cfg.HttpPort))
	httpListenAds := net.JoinHostPort(httpAds, httpPort)

	listener, err := net.Listen("tcp", httpListenAds)
	if err != nil {
		log.Fatalf("http server listen on address '%v' err: %v\n", httpListenAds, err)
	}
	defer listener.Close()

	// accept
	log.Printf("http server listen on address '%v'\n", httpListenAds)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept http request err: %v\n", err)
			continue
		}

		go handleHttpConnection(&conn)
	}
	log.Println("http server stop")
}

func handleHttpConnection(conn *net.Conn) {
	// init context
	ctx := &context{}
	defer func() {
		ctx.Close()

		if err := recover(); err != nil {
			log.Printf("unexpected err: %v\n", err)
		}
	}()

	// http handshake
	if err := httpHandshake(ctx, conn); err != nil {
		log.Printf("http handshake err: %v\n", err)
		return
	}

	// switch
	hikaricommon.Switch(ctx.localConn, ctx.serverConn, ctx.crypto, switchBufSize, switchTimeoutSec)
}

func httpHandshake(ctx *context, conn *net.Conn) error {
	ctx.localConn = conn

	// deadline
	deadline := time.Now().Add(time.Second * handshakeTimeoutSec)

	// set local connection timeout
	if err := (*conn).SetDeadline(deadline); err != nil {
		return err
	}

	// read http connect req
	var httpReq *httpConnectReq
	if req, err := readHttpReq(ctx.localConn); err != nil {
		return err
	} else {
		httpReq = req
	}

	// init crypto
	if c, err := hikaricommon.NewAESCrypto(secretKey, nil); err != nil {
		return err
	} else {
		crypto := interface{}(c).(hikaricommon.Crypto)
		ctx.crypto = &crypto
	}

	// send hikari req
	if serverConn, err := sendHikariReq(ctx.crypto, hikaricommon.HikariAddressTypeDomainName, httpReq.adsAndPort); err != nil {
		ctx.serverConn = serverConn
		return err
	} else {
		ctx.serverConn = serverConn

		// deadline
		if err := (*serverConn).SetDeadline(deadline); err != nil {
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

	// reply http
	if err := replyHttpReq(ctx.localConn, hikariRsp); err != nil {
		return err
	}

	return nil
}

func readHttpReq(conn *net.Conn) (*httpConnectReq, error) {
	var httpReq *http.Request
	if req, err := http.ReadRequest(bufio.NewReader(*conn)); err != nil {
		return nil, err
	} else {
		httpReq = req
	}

	// method
	if connectMethod != httpReq.Method {
		return nil, badHttpConnectReqErr
	}

	host := httpReq.Host
	var tgtAds []byte
	var tgtPort int

	if i := strings.LastIndex(host, ":"); i != -1 {
		tgtAds = []byte(host[:i])
		if p, err := strconv.Atoi(host[i+1:]); err != nil {
			return nil, err
		} else {
			tgtPort = p
		}

	} else {
		tgtAds = []byte(host)
		tgtPort = 80 // default port 80
	}

	adsLen := len(tgtAds)
	adsAndPort := make([]byte, adsLen+3)
	adsAndPort[0] = byte(adsLen)
	copy(adsAndPort[1:], tgtAds)
	binary.BigEndian.PutUint16(adsAndPort[adsLen+1:], uint16(tgtPort))

	return &httpConnectReq{adsAndPort}, nil
}

func replyHttpReq(conn *net.Conn, rsp *hikariRsp) error {
	c := *conn

	switch rsp.reply {
	case hikaricommon.HikariReplyOk:
		if _, err := c.Write(httpProxyOK); err != nil {
			return err
		}
		return nil

	case hikaricommon.HikariReplyVersionNotSupported:
		if _, err := c.Write(httpProxyFail); err != nil {
			return err
		}
		return hikaricommon.HikariVerNotSupportedErr

	case hikaricommon.HikariReplyAuthFail:
		if _, err := c.Write(httpProxyFail); err != nil {
			return err
		}
		return hikaricommon.HikariAuthFailErr

	case hikaricommon.HikariAdsTypeNotSupported:
		if _, err := c.Write(httpProxyFail); err != nil {
			return err
		}
		return hikaricommon.HikariAdsTypeNotSupportedErr

	case hikaricommon.HikariReplyDnsLookupFail:
		if _, err := c.Write(httpProxyFail); err != nil {
			return err
		}
		return hikaricommon.HikariDnsLookupFailErr

	case hikaricommon.HikariReplyConnectTargetFail:
		if _, err := c.Write(httpProxyFail); err != nil {
			return err
		}
		return hikaricommon.HikariConnectToTargetFailErr

	default:
		if _, err := c.Write(httpProxyFail); err != nil {
			return err
		}
		return hikaricommon.BadHikariRspErr
	}
}
