// thape - casting container images to gzipped tarballs.
// (c) 2025 Star Inc.

package config

import "github.com/star-inc/nui.go"

// init general config
var (
	AppMode = nui.String(envAppMode)
)

// init http config
var (
	HttpHost = nui.String(envHttpHost)
	HttpPort = nui.Integer(envHttpPort, nui.IntModeInt).(int)
)
