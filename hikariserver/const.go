package hikariserver

import "time"

const (
	// buffer size
	hikariReqBufSize = 276
	switchBufSize    = 4096

	// timeout
	handshakeTimeout time.Duration = 30
	dialTimeout      time.Duration = 30
	switchTimeout    time.Duration = 600
)
