package hikariclient

import "time"

const (
	// buffer size
	socksAuthReqBufSize = 257
	socksReqBufSize     = 262
	hikariRspBufSize    = 21
	switchBufSize       = 4096

	// timeout
	handshakeTimeoutSec time.Duration = 10
	dialTimeoutSec      time.Duration = 10
	switchTimeoutSec    time.Duration = 300
)
