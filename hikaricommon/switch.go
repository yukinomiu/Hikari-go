package hikaricommon

import (
	"net"
	"time"
)

func Switch(plainConn *net.Conn, encryptedConn *net.Conn, crypto *Crypto, bufSize int, timeoutSec time.Duration) {
	go pipePlain(plainConn, encryptedConn, crypto, bufSize, timeoutSec)
	pipeEncrypted(encryptedConn, plainConn, crypto, bufSize, timeoutSec)
}

func pipePlain(src *net.Conn, dst *net.Conn, crypto *Crypto, bufSize int, timeoutSec time.Duration) {
	s := *src
	d := *dst
	c := *crypto
	buf := make([]byte, bufSize)

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

func pipeEncrypted(src *net.Conn, dst *net.Conn, crypto *Crypto, bufSize int, timeoutSec time.Duration) {
	s := *src
	d := *dst
	c := *crypto
	buf := make([]byte, bufSize)

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
