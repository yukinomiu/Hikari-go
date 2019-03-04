package hikariclient

import "time"

const (
	// context
	ctxQueueSize = 32

	// buffer size
	ctxBufSize          = 4096
	socksAuthReqBufSize = 257
	socksReqBufSize     = 262
	hikariRspBufSize    = 21

	// timeout
	handshakeTimeout time.Duration = 30
	dialTimeout      time.Duration = 3
	switchTimeout    time.Duration = 600
)
