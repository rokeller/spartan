package server

import (
	"time"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port             uint16
	StaticContentDir string
	PathRoot         string

	// FallbackToIndex indicates whether spartan should respond with the
	// index.html of the configured root directory ([ServerConfig.StaticContentDir])
	// when a resource is not found. Defaults to true when set to nil.
	//
	// Deprecated: Use the settings from the [ServerConfig.NotFoundBehavior]
	// instead as it gives more control over the behavior when a requested
	// resource is not found.
	FallbackToIndex  *bool
	NotFoundBehavior *NotFoundBehavior

	Cache    Cache
	Security SecurityConfig
}

func (c *ServerConfig) GetNotFoundBehavior() *NotFoundBehavior {
	if c.FallbackToIndex == nil && c.NotFoundBehavior == nil {
		return &DefaultNotFoundBehavior
	} else if c.NotFoundBehavior != nil {
		return c.NotFoundBehavior
	} else {
		return &NotFoundBehavior{
			FallbackToIndex: c.FallbackToIndex,
		}
	}
}

type Cache struct {
	DefaultPolicy *CachePolicy
	Routes        []RouteMatchingCachePolicy
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

type RouteMatchingCachePolicy struct {
	Match  RouteMatcher
	Policy CachePolicy `mapstructure:",squash"`
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
