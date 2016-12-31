package lib

import "bytes"

// Base ...
type Base struct{}

// Name ...
func (p Base) Name() string {
	return "base"
}

// Extract ...
func (p Base) Extract(hostname string, bb *bytes.Buffer) (res map[string]string, err error) {
	return map[string]string{
		"hostname": hostname,
		"message":  bb.String(),
	}, nil
}
