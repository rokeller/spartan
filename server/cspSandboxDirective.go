package server

import (
	"context"
	"strings"
)

type SandboxDirectiveValue interface {
	DirectiveValue
	SandboxDirectiveValueMarker()
}

var _ SandboxDirectiveValue = SandboxAll{}
var _ SandboxDirectiveValue = SandboxWithAllowed{}

// SandboxAll represents the sandbox directive that applies all sandbox restrictions.
type SandboxAll struct{}

func (SandboxAll) Value(ctx context.Context) string {
	return ""
}
func (SandboxAll) SandboxDirectiveValueMarker() {}

// SandboxWithAllowed represents allowed (i.e., removed) sandbox restrictions.
type SandboxWithAllowed []SandboxAllow

func (s SandboxWithAllowed) Value(ctx context.Context) string {
	if len(s) > 0 {
		v := make([]string, len(s))
		for i, a := range s {
			v[i] = string(a)
		}
		return strings.Join(v, " ")
	}
	return ""
}
func (SandboxWithAllowed) SandboxDirectiveValueMarker() {}

type SandboxAllow string
