package hikariserver

import (
	"crypto/md5"
	"encoding/hex"
)

var (
	authMap   map[string]int
	secretKey []byte
)

func initStatus() {
	// init auth
	authMap = make(map[string]int, len(cfg.PrivateKeyList))

	for _, k := range cfg.PrivateKeyList {
		authArray := md5.Sum([]byte(k))
		s := hex.EncodeToString(authArray[:])
		authMap[s] = 0
	}

	// init secret key
	secretKeyArray := md5.Sum([]byte(cfg.Secret))
	secretKey = secretKeyArray[:]
}

func isValidAuth(authHex string) bool {
	_, exits := authMap[authHex]
	return exits
}
