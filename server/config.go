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
	PermissionsPolicy             *PermissionsPolicy
	ReferrerPolicy                *ReferrerPolicy
	ReportingEndpoints            *ReportingEndpoints
	StrictTransportSecurityPolicy *StrictTransportSecurityPolicy
}

func (c *SecurityConfig) GetContentSecurityPolicy() *ContentSecurityPolicy {
	csp := c.ContentSecurityPolicy
	if nil == csp {
		return &DefaultContentSecurityPolicy
	}
	return csp
}

func (c *SecurityConfig) GetContentTypeOptionsNoSniff() bool {
	if nil == c.ContentTypeOptionsNoSniff {
		c.ContentTypeOptionsNoSniff = new(bool)
		*c.ContentTypeOptionsNoSniff = true
	}
	return *c.ContentTypeOptionsNoSniff
}

func (c *SecurityConfig) GetPermissionsPolicy() *PermissionsPolicy {
	pp := c.PermissionsPolicy
	if nil == pp {
		return &DefaultPermissionsPolicy
	}
	return pp
}

func (c *SecurityConfig) GetReferrerPolicy() *ReferrerPolicy {
	if nil == c.ReferrerPolicy {
		return &DefaultReferrerPolicy
	}
	return c.ReferrerPolicy
}

func (c *SecurityConfig) GetStrictTransportSecurityPolicy() *StrictTransportSecurityPolicy {
	if nil == c.StrictTransportSecurityPolicy {
		return &DefaultStrictTransportSecurityPolicy
	}
	return c.StrictTransportSecurityPolicy
}
