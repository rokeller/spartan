package server

import (
	"fmt"
	"strings"
)

type ReportingEndpoint struct {
	Name string
	Url  string
}

type ReportingEndpoints []ReportingEndpoint

func (e ReportingEndpoints) ReportingEndpointsHeaderValue() string {
	directives := make([]string, len(e))

	for i, ep := range e {
		directives[i] = fmt.Sprintf("%s=%q", ep.Name, ep.Url)
	}

	return strings.Join(directives, ", ")
}

func (e ReportingEndpoints) HasEndpoints() bool {
	return len(e) > 0
}

func (e ReportingEndpoints) GetEndpoint(name string) *ReportingEndpoint {
	for _, ep := range e {
		if ep.Name == name {
			return &ep
		}
	}
	return nil
}
