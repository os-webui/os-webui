package sdk

import (
	"github.com/os-webui/os-webui/version"
)

func Platform() string {
	return version.Platform
}
func Version() string {
	return version.Version
}
func Date() string {
	return version.Date
}
func Commit() string {
	return version.Commit
}