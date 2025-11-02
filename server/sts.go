package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	DefaultStrictTransportSecurityPolicy = StrictTransportSecurityPolicy{
		IncludeSubDomains: true,
		MaxAge:            time.Hour * 24 * 365,
	}
)

type StrictTransportSecurityPolicy struct {
	IncludeSubDomains bool
	MaxAge            time.Duration
}

func (p StrictTransportSecurityPolicy) AddToResponse(w http.ResponseWriter) {
	w.Header().Add("strict-transport-security", p.StsHeaderValue())
}

func (p StrictTransportSecurityPolicy) StsHeaderValue() string {
	directives := []string{
		fmt.Sprintf("max-age=%d",
			int64(p.MaxAge.Truncate(time.Second).Seconds())),
	}

	if p.IncludeSubDomains {
		directives = append(directives, "includeSubDomains")
	}

	return strings.Join(directives, "; ")
}
