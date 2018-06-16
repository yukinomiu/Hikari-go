package hikariclient

import (
	"hikari-go/hikaricommon"
	"io"
	"net"
	"time"
)

type hikariRsp struct {
	reply          byte
	bindAdsType    byte
	bindAdsAndPort []byte
	extraData      []byte
}

func sendHikariReq(crypto *hikaricommon.Crypto, hikariAdsType byte, adsAndPort []byte) (*net.Conn, error) {
	cr := *crypto

	serverConn, err := net.DialTimeout("tcp", serverAds, time.Second*dialTimeoutSec)
	if err != nil {
		return nil, err
	}

	// send iv
	if _, err := serverConn.Write(cr.GetIV()); err != nil {
		return &serverConn, err
	}

	// send hikari request
	buf := make([]byte, 18+len(adsAndPort))
	buf[0] = hikaricommon.HikariVer1
	copy(buf[1:17], auth)
	buf[17] = hikariAdsType
	copy(buf[18:], adsAndPort)
	cr.Encrypt(buf)

	if _, err := serverConn.Write(buf); err != nil {
		return &serverConn, err
	}

	return &serverConn, nil
}

func readHikariRsp(crypto *hikaricommon.Crypto, conn *net.Conn) (*hikariRsp, error) {
	c := *conn
	cr := *crypto

	// init buffer
	buf := make([]byte, hikariRspBufSize)

	// read
	n, err := io.ReadAtLeast(c, buf, 2)
	if err != nil {
		return nil, err
	}
	cr.Decrypt(buf[:n])

	// ver
	ver := buf[0]
	if ver != hikaricommon.HikariVer1 {
		return nil, hikaricommon.HikariVerNotSupportedErr
	}

	// reply
	reply := buf[1]

	switch reply {
	case hikaricommon.HikariReplyOk:
		if n < 4 {
			b := buf[n:]
			r, err := io.ReadAtLeast(c, b, 4-n)
			if err != nil {
				return nil, err
			}
			cr.Decrypt(b[:r])
			n += r
		}

		// bind address type
		bindAdsType := buf[2]
		var expectLen int
		var extraData []byte

		switch bindAdsType {
		case hikaricommon.HikariAddressTypeIpv4:
			expectLen = 5 + net.IPv4len

		case hikaricommon.HikariAddressTypeIpv6:
			expectLen = 5 + net.IPv6len

		default:
			return nil, hikaricommon.HikariAdsTypeNotSupportedErr
		}

		if n == expectLen {
		} else if n < expectLen {
			if expectLen > hikariRspBufSize {
				return nil, hikaricommon.BadHikariRspErr
			}

			b := buf[n:expectLen]
			if _, err := io.ReadFull(c, b); err != nil {
				return nil, err
			}
			cr.Decrypt(b)
		} else if n > expectLen {
			extraData = buf[expectLen:n]
			cr.Decrypt(extraData)
		}

		bindAdsAndPort := buf[3:expectLen]
		return &hikariRsp{reply, bindAdsType, bindAdsAndPort, extraData}, nil

	case hikaricommon.HikariReplyVersionNotSupported:
		fallthrough

	case hikaricommon.HikariReplyAuthFail:
		fallthrough

	case hikaricommon.HikariReplyDnsLookupFail:
		fallthrough

	case hikaricommon.HikariReplyConnectTargetFail:
		return &hikariRsp{reply: reply}, nil

	default:
		return nil, hikaricommon.BadHikariRspErr
	}
}
