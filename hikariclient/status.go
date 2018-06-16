package hikariclient

import (
	"crypto/md5"
	"net"
	"strconv"
)

var (
	auth      []byte
	serverAds string
	secretKey []byte
)

func initStatus() {
	// init auth
	authArray := md5.Sum([]byte(cfg.PrivateKey))
	auth = authArray[:]

	// init server address
	srvAds := cfg.ServerAddress
	srvPort := strconv.Itoa(int(cfg.ServerPort))
	serverAds = net.JoinHostPort(srvAds, srvPort)

	// init secret key
	secretKeyArray := md5.Sum([]byte(cfg.Secret))
	secretKey = secretKeyArray[:]
}
