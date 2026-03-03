// Package updater provides version manifest management for dnstm binaries.
package updater

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const (
	// VersionManifestFile is the filename for the version manifest.
	VersionManifestFile = "versions.json"
)

// VersionManifest stores installed versions of transport binaries.
type VersionManifest struct {
	SlipstreamServer string    `json:"slipstream-server,omitempty"`
	SSServer         string    `json:"ssserver,omitempty"`
	Microsocks       string    `json:"microsocks,omitempty"`
	SSHTunUser       string    `json:"sshtun-user,omitempty"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// GetManifestPath returns the path to the version manifest file.
func GetManifestPath() string {
	return filepath.Join("/etc/dnstm", VersionManifestFile)
}

// LoadManifest loads the version manifest from disk.
func LoadManifest() (*VersionManifest, error) {
	path := GetManifestPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &VersionManifest{}, nil
		}
		return nil, err
	}

	var manifest VersionManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// Save saves the version manifest to disk.
func (m *VersionManifest) Save() error {
	path := GetManifestPath()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	m.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GetVersion returns the installed version for a binary.
func (m *VersionManifest) GetVersion(binaryName string) string {
	switch binaryName {
	case "slipstream-server":
		return m.SlipstreamServer
	case "ssserver":
		return m.SSServer
	case "microsocks":
		return m.Microsocks
	case "sshtun-user":
		return m.SSHTunUser
	default:
		return ""
	}
}

// SetVersion sets the installed version for a binary.
func (m *VersionManifest) SetVersion(binaryName, version string) {
	switch binaryName {
	case "slipstream-server":
		m.SlipstreamServer = version
	case "ssserver":
		m.SSServer = version
	case "microsocks":
		m.Microsocks = version
	case "sshtun-user":
		m.SSHTunUser = version
	}
}
