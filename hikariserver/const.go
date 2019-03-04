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
	dialTimeout      time.Duration = 3
	switchTimeout    time.Duration = 600

	// connect count
	maxConnectCount        = 6
	maxConnectCountWithIp4 = 4
)
