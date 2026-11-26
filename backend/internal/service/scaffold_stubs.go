package service

import (
	"strings"

	infraerrors "github.com/x-frame/x-frame/internal/pkg/errors"
)

var (
	ErrRefreshTokenNotFound = infraerrors.Unauthorized("REFRESH_TOKEN_NOT_FOUND", "refresh token not found")
	ErrTurnstileNotConfigured = infraerrors.BadRequest("TURNSTILE_NOT_CONFIGURED", "turnstile not configured")
)

// MaxValidityDays caps default-subscription validity in settings parsing.
const MaxValidityDays = 3650

// OpsQueryMode is retained for settings JSON compatibility.
type OpsQueryMode string

const (
	OpsQueryModeAuto OpsQueryMode = "auto"
)

// ParseOpsQueryMode normalizes ops query mode (ops domain removed; always auto).
func ParseOpsQueryMode(raw string) OpsQueryMode {
	_ = raw
	return OpsQueryModeAuto
}

// CompareVersions compares two dotted numeric version strings.
// Returns <0 if a<b, 0 if equal, >0 if a>b.
func CompareVersions(a, b string) int {
	as := strings.Split(a, ".")
	bs := strings.Split(b, ".")
	n := len(as)
	if len(bs) > n {
		n = len(bs)
	}
	for i := 0; i < n; i++ {
		var ai, bi int
		if i < len(as) {
			ai = atoiSafe(as[i])
		}
		if i < len(bs) {
			bi = atoiSafe(bs[i])
		}
		if ai != bi {
			return ai - bi
		}
	}
	return 0
}

func atoiSafe(s string) int {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			break
		}
		n = n*10 + int(r-'0')
	}
	return n
}
