package hikariserver

import "time"

const (
	// buffer size
	hikariReqBufSize = 276
	switchBufSize    = 4096

	// timeout
	handshakeTimeoutSec time.Duration = 30
	dialTimeoutSec      time.Duration = 30
	switchTimeoutSec    time.Duration = 600
)
