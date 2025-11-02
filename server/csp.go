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

	// Document directives
	BaseUri NoneOrSourceExpressionList
	Sandbox SandboxDirectiveValue

	// Navigation directives
	FormAction, FrameAncestors NoneOrSourceExpressionList

	// Reporting directives
	ReportTo string

}

func (p ContentSecurityPolicy) HeaderValue(ctx context.Context) string {
	directives := []string{}

	// Fetch directives
	directives = appendDirectives(directives, ctx, p.ChildSrc, "child-src")
	directives = appendDirectives(directives, ctx, p.ConnectSrc, "connect-src")
	directives = appendDirectives(directives, ctx, p.DefaultSrc, "default-src")
	directives = appendDirectives(directives, ctx, p.FencedFrameSrc, "fenced-frame-src")
	directives = appendDirectives(directives, ctx, p.FontSrc, "font-src")
	directives = appendDirectives(directives, ctx, p.FrameSrc, "frame-src")
	directives = appendDirectives(directives, ctx, p.ImgSrc, "img-src")
	directives = appendDirectives(directives, ctx, p.ManifestSrc, "manifest-src")
	directives = appendDirectives(directives, ctx, p.MediaSrc, "media-src")
	directives = appendDirectives(directives, ctx, p.ObjectSrc, "object-src")
	directives = appendDirectives(directives, ctx, p.ScriptSrc, "script-src")
	directives = appendDirectives(directives, ctx, p.ScriptSrcElem, "script-src-elem")
	directives = appendDirectives(directives, ctx, p.ScriptSrcElem, "script-src-attr")
	directives = appendDirectives(directives, ctx, p.StyleSrc, "style-src")
	directives = appendDirectives(directives, ctx, p.StyleSrcElem, "style-src-elem")
	directives = appendDirectives(directives, ctx, p.StyleSrcAttr, "style-src-attr")
	directives = appendDirectives(directives, ctx, p.WorkerSrc, "worker-src-src")

	// Document directives
	directives = appendDirectives(directives, ctx, p.BaseUri, "base-uri")
	directives = appendDirectives(directives, ctx, p.Sandbox, "sandbox")

	// Navigation directives
	directives = appendDirectives(directives, ctx, p.FormAction, "form-action")
	directives = appendDirectives(directives, ctx, p.FrameAncestors, "frame-ancestors")

	// Reporting directives
	if p.ReportTo != "" {
		directives = appendDirectivesValue(directives, "report-to", p.ReportTo)
	}

	if len(directives) > 0 {
		return strings.Join(directives, "; ")
	}
	return ""
}

func appendDirectives(directives []string, ctx context.Context, v DirectiveValue, name string) []string {
	if nil != v {
		return appendDirectivesValue(directives, name, v.Value(ctx))
	}

	return directives
}

func appendDirectivesValue(directives []string, name, value string) []string {
	var directive string
	if value != "" {
		directive = fmt.Sprintf("%s %s", name, value)
	} else {
		directive = name
	}
	return append(directives, directive)
}
