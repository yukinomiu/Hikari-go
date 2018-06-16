package hikaricommon

import "errors"

const (
	// version
	HikariVer1 byte = 1

	// address type
	HikariAddressTypeIpv4       byte = 0
	HikariAddressTypeIpv6       byte = 1
	HikariAddressTypeDomainName byte = 2

	// reply
	HikariReplyOk                  byte = 0
	HikariReplyVersionNotSupported byte = 1
	HikariReplyAuthFail            byte = 2
	HikariAdsTypeNotSupported      byte = 3
	HikariReplyDnsLookupFail       byte = 4
	HikariReplyConnectTargetFail   byte = 5
)

var (
	// hikari errors
	HikariVerNotSupportedErr     = errors.New("hikari version not supported")
	HikariAuthFailErr            = errors.New("hikari auth fail")
	HikariAdsTypeNotSupportedErr = errors.New("hikari address type not supported")
	HikariDnsLookupFailErr       = errors.New("dns lookup fail")
	HikariConnectToTargetFailErr = errors.New("connect to target fail")

	// other
	BadHikariRspErr = errors.New("bad hikari response")
	BadHikariReqErr = errors.New("bad hikari request")
)
