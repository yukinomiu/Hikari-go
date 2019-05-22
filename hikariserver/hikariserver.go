package hikariserver

import (
	"crypto/aes"
	"encoding/binary"
	"encoding/hex"
	"hikari-go/hikaricommon"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type hikariReq struct {
	adsType byte
	ads     []byte
	port    []byte
}

func startHikariServer() {
	if cfg.ListenPort == 0 {
		log.Println("hikari server disabled")
		return
	}

	// init listener
	ads := cfg.ListenAddress
	port := strconv.Itoa(int(cfg.ListenPort))
	listenAds := net.JoinHostPort(ads, port)

	listener, err := net.Listen("tcp", listenAds)
	if err != nil {
		log.Fatalf("hikari server listen on address '%v' err: %v\n", listenAds, err)
	}
	defer func() {
		_ = listener.Close()
	}()

	// accept
	log.Printf("hikari server listen on address '%v'\n", listenAds)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept hikari request err: %v\n", err)
			continue
		}

		go handleConnection(&conn)
	}
	// log.Println("hikari server stop")
}

func handleConnection(conn *net.Conn) {
	// init context
	ctx := getContext()
	defer func() {
		ctx.Close()
		returnContext(ctx)

		if err := recover(); err != nil {
			log.Printf("unexpected err: %v\n", err)
		}
	}()

	// hikari handshake
	if err := hikariHandshake(ctx, conn); err != nil {
		log.Printf("hikari handshake err: %v\n", err)
		return
	}

	// switch
	hikariCtx := interface{}(ctx).(hikaricommon.Context)
	hikaricommon.Switch(&hikariCtx, switchTimeout)
}

func hikariHandshake(ctx *context, conn *net.Conn) error {
	ctx.clientConn = conn

	// deadline
	deadline := time.Now().Add(time.Second * handshakeTimeout)

	// set client connection timeout
	if err := (*conn).SetDeadline(deadline); err != nil {
		return err
	}

	// read IV
	var iv []byte
	if i, err := readIV(ctx.clientConn); err != nil {
		return err
	} else {
		iv = i
	}

	// init crypto
	if c, err := hikaricommon.NewAESCrypto(secretKey, iv); err != nil {
		return err
	} else {
		crypto := interface{}(c).(hikaricommon.Crypto)
		ctx.crypto = &crypto
	}

	// read hikari req
	var hikariReq *hikariReq
	if req, err := readHikariReq(ctx.clientConn, ctx.crypto); err != nil {
		return err
	} else {
		hikariReq = req
	}

	// connect target
	if c, err := connectTarget(ctx.clientConn, hikariReq, ctx.crypto); err != nil {
		ctx.targetConn = c
		return err
	} else {
		ctx.targetConn = c

		// deadline
		if err := (*c).SetDeadline(deadline); err != nil {
			return err
		}
	}

	// reply hikari req
	if err := replyHikariReq(ctx.clientConn, ctx.targetConn, ctx.crypto); err != nil {
		return err
	}

	return nil
}

func readIV(conn *net.Conn) ([]byte, error) {
	c := *conn

	// init buffer
	buf := make([]byte, aes.BlockSize)

	// read
	if _, err := io.ReadFull(c, buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func readHikariReq(conn *net.Conn, crypto *hikaricommon.Crypto) (*hikariReq, error) {
	c := *conn
	cr := *crypto

	// init buffer
	buf := make([]byte, hikariReqBufSize)

	// read
	n, err := io.ReadAtLeast(c, buf, 19)
	if err != nil {
		return nil, err
	}
	cr.Decrypt(buf[:n])

	// ver
	ver := buf[0]
	if ver != hikaricommon.HikariVer1 {
		rsp := []byte{hikaricommon.HikariVer1, hikaricommon.HikariReplyVersionNotSupported}
		cr.Encrypt(rsp)
		if _, err := c.Write(rsp); err != nil {
			return nil, err
		}

		return nil, hikaricommon.HikariVerNotSupportedErr
	}

	// auth
	authHex := hex.EncodeToString(buf[1:17])
	if !isValidAuth(authHex) {
		rsp := []byte{hikaricommon.HikariVer1, hikaricommon.HikariReplyAuthFail}
		cr.Encrypt(rsp)
		if _, err := c.Write(rsp); err != nil {
			return nil, err
		}

		return nil, hikaricommon.HikariAuthFailErr
	}

	// address type
	adsType := buf[17]
	var expectLen int
	var ads, port []byte

	switch adsType {
	case hikaricommon.HikariAddressTypeDomainName:
		adsLen := int(buf[18])
		expectLen = 21 + adsLen

		i := 19 + adsLen
		ads = buf[19:i]
		port = buf[i:expectLen]

	case hikaricommon.HikariAddressTypeIpv4:
		expectLen = 20 + net.IPv4len

		i := 18 + net.IPv4len
		ads = buf[18:i]
		port = buf[i:expectLen]

	case hikaricommon.HikariAddressTypeIpv6:
		expectLen = 20 + net.IPv6len

		i := 18 + net.IPv6len
		ads = buf[18:i]
		port = buf[i:expectLen]

	default:
		rsp := []byte{hikaricommon.HikariVer1, hikaricommon.HikariAdsTypeNotSupported}
		cr.Encrypt(rsp)
		if _, err := c.Write(rsp); err != nil {
			return nil, err
		}

		return nil, hikaricommon.HikariAdsTypeNotSupportedErr
	}

	if n == expectLen {
	} else if n < expectLen {
		if expectLen > hikariReqBufSize {
			return nil, hikaricommon.BadHikariReqErr
		}

		b := buf[n:expectLen]
		if _, err := io.ReadFull(c, b); err != nil {
			return nil, err
		}
		cr.Decrypt(b)
	} else if n > expectLen {
		return nil, hikaricommon.BadHikariReqErr
	}

	return &hikariReq{adsType, ads, port}, nil
}

func connectTarget(clientConn *net.Conn, hikariReq *hikariReq, crypto *hikaricommon.Crypto) (*net.Conn, error) {
	c := *clientConn
	cr := *crypto

	var ips []net.IP
	if hikariReq.adsType == hikaricommon.HikariAddressTypeDomainName {
		// DNS lookup
		host := string(hikariReq.ads)
		if ipSlice, err := net.LookupIP(host); err != nil {
			log.Printf("dns lookup '%v' err: %v\n", host, err)

			rsp := []byte{hikaricommon.HikariVer1, hikaricommon.HikariReplyDnsLookupFail}
			cr.Encrypt(rsp)
			if _, err := c.Write(rsp); err != nil {
				return nil, err
			}

			return nil, hikaricommon.HikariDnsLookupFailErr
		} else {
			ips = ipSlice
		}

	} else {
		ips = []net.IP{hikariReq.ads}
	}

	var tgtConn net.Conn

	portStr := strconv.Itoa(int(binary.BigEndian.Uint16(hikariReq.port)))
	hasIp4 := false
	for i, ip := range ips {
		if !hasIp4 && ip.To4() != nil {
			hasIp4 = true
		}

		if i >= maxConnectCount || (hasIp4 && i >= maxConnectCountWithIp4) {
			break
		}

		tgtAdsStr := net.JoinHostPort(ip.String(), portStr)
		if c, err := net.DialTimeout("tcp", tgtAdsStr, time.Second*dialTimeout); err != nil {
			log.Printf("connect target '%v' (try count %v) err: %v\n", tgtAdsStr, i, err)
			continue
		} else {
			tgtConn = c
			break
		}
	}

	if tgtConn == nil {
		rsp := []byte{hikaricommon.HikariVer1, hikaricommon.HikariReplyConnectTargetFail}
		cr.Encrypt(rsp)
		if _, err := c.Write(rsp); err != nil {
			return nil, err
		}

		return nil, hikaricommon.HikariConnectToTargetFailErr
	}

	return &tgtConn, nil
}

func replyHikariReq(clientConn *net.Conn, targetConn *net.Conn, crypto *hikaricommon.Crypto) error {
	c := *clientConn
	cr := *crypto

	bindAdsType, bindAds, bindPort := getBindInfo(targetConn)
	bindAdsLen := len(bindAds)
	i := 3 + bindAdsLen

	buf := make([]byte, 5+bindAdsLen)
	buf[0] = hikaricommon.HikariVer1
	buf[1] = hikaricommon.HikariReplyOk
	buf[2] = bindAdsType
	copy(buf[3:i], bindAds)
	binary.BigEndian.PutUint16(buf[i:], bindPort)
	cr.Encrypt(buf)

	if _, err := c.Write(buf); err != nil {
		return err
	}

	return nil
}

func getBindInfo(conn *net.Conn) (byte, []byte, uint16) {
	localAddrStr := (*conn).LocalAddr().String()
	i := strings.LastIndex(localAddrStr, ":")

	// ip
	var bindAdsType byte
	var bindAds []byte
	localIp := net.ParseIP(localAddrStr[:i])
	if ip4 := localIp.To4(); ip4 != nil {
		// v4
		bindAdsType = hikaricommon.HikariAddressTypeIpv4
		bindAds = ip4
	} else {
		// v6
		bindAdsType = hikaricommon.HikariAddressTypeIpv6
		bindAds = localIp.To16()
	}

	// port
	bindPort, _ := strconv.Atoi(localAddrStr[i+1:])

	return bindAdsType, bindAds, uint16(bindPort)
}
