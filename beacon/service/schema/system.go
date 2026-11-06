// Package schema
// Date: 2024/3/6 11:21
// Author: Amu
// Description:
package schema

// SystemVersion is the wire representation of the running Beacon Server build
// and its latest update-check result.
//
// The Server never upgrades itself: the version endpoint only reports state so
// the console can tell operators that a newer Beacon release exists.
type SystemVersion struct {
	// Version is the current build version of the running Beacon Server.
	Version string `json:"version"`
	// UpdateAvailable reports whether a newer Beacon release is known.
	UpdateAvailable bool `json:"update_available"`
	// LatestVersion is the newest known release version; empty when unknown.
	LatestVersion string `json:"latest_version,omitempty"`
	// ReleaseNotes is the release note of the newest known release.
	ReleaseNotes string `json:"release_notes,omitempty"`
	// CheckedAt is the RFC3339 timestamp of the last successful update check.
	CheckedAt string `json:"checked_at,omitempty"`
}
