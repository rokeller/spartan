package server

import (
	"context"
	"fmt"
	"strings"
)

var (
	DefaultContentSecurityPolicy = ContentSecurityPolicy{
		ReportOnly:     false,
		DefaultSrc:     SelfDirectiveValue{},
		FormAction:     SelfDirectiveValue{},
		FrameAncestors: SelfDirectiveValue{},
	}
)

type ContentSecurityPolicy struct {
	ReportOnly bool

	// Fetch directives
	ChildSrc, ConnectSrc, DefaultSrc, FencedFrameSrc, FontSrc, FrameSrc, ImgSrc,
	ManifestSrc, MediaSrc, ObjectSrc, ScriptSrc, ScriptSrcElem, ScriptSrcAttr,
	StyleSrc, StyleSrcElem, StyleSrcAttr, WorkerSrc FetchDirectiveValue

	// Navigation directives
	FormAction, FrameAncestors NoneOrSourceExpressionList
}

func (p ContentSecurityPolicy) HeaderValue(ctx context.Context) string {
	directives := []string{}

	directives = AppendDirectives(directives, ctx, p.ChildSrc, "child-src")
	directives = AppendDirectives(directives, ctx, p.ConnectSrc, "connect-src")
	directives = AppendDirectives(directives, ctx, p.DefaultSrc, "default-src")
	directives = AppendDirectives(directives, ctx, p.FencedFrameSrc, "fenced-frame-src")
	directives = AppendDirectives(directives, ctx, p.FontSrc, "font-src")
	directives = AppendDirectives(directives, ctx, p.FrameSrc, "frame-src")
	directives = AppendDirectives(directives, ctx, p.ImgSrc, "img-src")
	directives = AppendDirectives(directives, ctx, p.ManifestSrc, "manifest-src")
	directives = AppendDirectives(directives, ctx, p.MediaSrc, "media-src")
	directives = AppendDirectives(directives, ctx, p.ObjectSrc, "object-src")
	directives = AppendDirectives(directives, ctx, p.ScriptSrc, "script-src")
	directives = AppendDirectives(directives, ctx, p.ScriptSrcElem, "script-src-elem")
	directives = AppendDirectives(directives, ctx, p.ScriptSrcElem, "script-src-attr")
	directives = AppendDirectives(directives, ctx, p.StyleSrc, "style-src")
	directives = AppendDirectives(directives, ctx, p.StyleSrcElem, "style-src-elem")
	directives = AppendDirectives(directives, ctx, p.StyleSrcAttr, "style-src-attr")
	directives = AppendDirectives(directives, ctx, p.WorkerSrc, "worker-src-src")

	directives = AppendDirectives(directives, ctx, p.FormAction, "form-action")
	directives = AppendDirectives(directives, ctx, p.FrameAncestors, "frame-ancestors")

	if len(directives) > 0 {
		return strings.Join(directives, "; ")
	}
	return ""
}

func AppendDirectives(directives []string, ctx context.Context, v DirectiveValue, name string) []string {
	if nil != v {
		directive := fmt.Sprintf("%s %s", name, v.Value(ctx))
		return append(directives, directive)
	}

	return directives
}
