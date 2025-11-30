package server

import "net/http"

var (
	DefaultReferrerPolicy = ReferrerPolicy{
		Value: SameOriginReferrerPolicyValue,
	}
)

type ReferrerPolicy struct {
	Value ReferrerPolicyValue
}

func (p ReferrerPolicy) AddToResponse(w http.ResponseWriter) {
	w.Header().Add("referrer-policy", p.Value.Value())
}

type ReferrerPolicyValue interface {
	Value() string
}

type referrerPolicyStringValue string

var _ ReferrerPolicyValue = referrerPolicyStringValue("")

func (r referrerPolicyStringValue) Value() string {
	return string(r)
}

const (
	NoReferrerReferrerPolicyValue                  = referrerPolicyStringValue("no-referrer")
	NoReferrerWhenDowngradeReferrerPolicyValue     = referrerPolicyStringValue("no-referrer-when-downgrade")
	OriginReferrerPolicyValue                      = referrerPolicyStringValue("origin")
	OriginWhenCrossOriginReferrerPolicyValue       = referrerPolicyStringValue("origin-when-cross-origin")
	SameOriginReferrerPolicyValue                  = referrerPolicyStringValue("same-origin")
	StrictOriginReferrerPolicyValue                = referrerPolicyStringValue("strict-origin")
	StrictOriginWhenCrossOriginReferrerPolicyValue = referrerPolicyStringValue("strict-origin-when-cross-origin")
	UnsafeUrlReferrerPolicyValue                   = referrerPolicyStringValue("unsafe-url")
)
