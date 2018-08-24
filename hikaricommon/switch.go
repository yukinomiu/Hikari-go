package hikaricommon

import (
	"net"
	"time"
)

func Switch(ctx *Context, timeoutSec time.Duration) {
	c := *ctx
	go pipePlain(c.GetPlainConn(), c.GetEncryptedConn(), c.GetCrypto(), c.GetPlainBuf(), timeoutSec)
	pipeEncrypted(c.GetEncryptedConn(), c.GetPlainConn(), c.GetCrypto(), c.GetEncryptedBuf(), timeoutSec)
}

func pipePlain(src *net.Conn, dst *net.Conn, crypto *Crypto, buf []byte, timeoutSec time.Duration) {
	s := *src
	d := *dst
	c := *crypto

	defer func() {
		s.Close()
		d.Close()
	}()

	var data []byte
	addTime := time.Second * timeoutSec

	for {
		// set timeout
		t := time.Now().Add(addTime)
		if err := s.SetReadDeadline(t); err != nil {
			break
		}
		if err := d.SetWriteDeadline(t); err != nil {
			break
		}

		// pipe
		if n, err := s.Read(buf); err != nil {
			if n != 0 {
				data = buf[:n]
				c.Encrypt(data)

				_, err = d.Write(data)
				if err != nil {
					break
				}
			}

			break

		} else {
			data = buf[:n]
			c.Encrypt(data)

			_, err = d.Write(data)
			if err != nil {
				break
			}
		}
	}
}

func pipeEncrypted(src *net.Conn, dst *net.Conn, crypto *Crypto, buf []byte, timeoutSec time.Duration) {
	s := *src
	d := *dst
	c := *crypto

	defer func() {
		s.Close()
		d.Close()
	}()

	var data []byte
	addTime := time.Second * timeoutSec

	for {
		// set timeout
		t := time.Now().Add(addTime)
		if err := s.SetReadDeadline(t); err != nil {
			break
		}
		if err := d.SetWriteDeadline(t); err != nil {
			break
		}

		// pipe
		if n, err := s.Read(buf); err != nil {
			if n != 0 {
				data = buf[:n]
				c.Decrypt(data)

				_, err = d.Write(data)
				if err != nil {
					break
				}
			}

			break

		} else {
			data = buf[:n]
			c.Decrypt(data)

			_, err = d.Write(data)
			if err != nil {
				break
			}
		}
	}
}
