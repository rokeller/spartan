package server

import "time"

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port             uint16
	StaticContentDir string
	PathRoot         string

	Cache    Cache
	Security SecurityConfig
}

type Cache struct {
	DefaultPolicy CachePolicy
}

type CachePolicy struct {
	Immutable       bool
	MustRevalidate  bool
	MustUnderstand  bool
	NoCache         bool
	NoStore         bool
	NoTransform     bool
	Private         bool
	ProxyRevalidate bool
	Public          bool

	MaxAge               *time.Duration
	SharedMaxAge         *time.Duration
	StaleIfError         *time.Duration
	StaleWhileRevalidate *time.Duration
}

type SecurityConfig struct {
	ContentTypeOptionsNoSniff     *bool
	ContentSecurityPolicy         *ContentSecurityPolicy
	StrictTransportSecurityPolicy *StrictTransportSecurityPolicy
}

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
