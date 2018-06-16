package hikariclient

import "errors"

const (
	// method
	connectMethod = "CONNECT"
)

var (
	// http errors
	badHttpConnectReqErr = errors.New("bad http connect request")
)
