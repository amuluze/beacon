// Package service
// Date: 2025/01/15
// Description: Version information for the Collia agent.
package service

// Version holds the build version of the Collia agent.
type Version struct {
	Version string
}

// NewVersion creates a Version value from the build version string.
func NewVersion(version string) Version {
	return Version{Version: version}
}

// String returns the version string.
func (v Version) String() string {
	return v.Version
}
