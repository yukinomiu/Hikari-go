package hikariserver

import (
	"time"
)

const (
	// context
	ctxQueueSize = 64

	// buffer size
	ctxBufSize       = 4096
	hikariReqBufSize = 276

	// timeout
	handshakeTimeout time.Duration = 30
	dialTimeout      time.Duration = 10
	switchTimeout    time.Duration = 600
)
