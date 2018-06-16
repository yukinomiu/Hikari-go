package hikariserver

import "time"

const (
	// buffer size
	hikariReqBufSize = 276
	switchBufSize    = 4096

	// timeout
	handshakeTimeoutSec time.Duration = 10
	dialTimeoutSec      time.Duration = 10
	switchTimeoutSec    time.Duration = 300
)
