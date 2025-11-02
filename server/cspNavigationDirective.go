package server

import (
	"context"
	"strings"
)

type NoneOrSourceExpressionList interface {
	DirectiveValue
	NoneOrSourceExpressionListMarker()
}

type SourceExpressionListItem interface {
	DirectiveValue
	SourceExpressionListItemMarker()
}

// SourceExpressionList combines multiple source expression list item values.
type SourceExpressionList []SourceExpressionListItem

func (l SourceExpressionList) Value(ctx context.Context) string {
	parts := make([]string, len(l))
	for i, val := range l {
		parts[i] = val.Value(ctx)
	}

	return strings.Join(parts, " ")
}
func (SourceExpressionList) NoneOrSourceExpressionListMarker() {}

var _ NoneOrSourceExpressionList = NoneDirectiveValue{}
var _ NoneOrSourceExpressionList = SourceExpressionList{}
var _ NoneOrSourceExpressionList = SelfDirectiveValue{}
var _ NoneOrSourceExpressionList = HostSourceDirectiveValue{}
var _ NoneOrSourceExpressionList = SchemeSourceDirectiveValue("")

var _ SourceExpressionListItem = SelfDirectiveValue{}
var _ SourceExpressionListItem = HostSourceDirectiveValue{}
var _ SourceExpressionListItem = SchemeSourceDirectiveValue("")
