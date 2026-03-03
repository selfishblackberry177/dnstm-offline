// Package binary provides binary management for external tools.
package binary

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// BinaryType identifies a binary.
type BinaryType string

const (
	// Server binaries (used in production)
	BinaryDNSTTServer      BinaryType = "dnstt-server"
	BinarySlipstreamServer BinaryType = "slipstream-server"
	BinarySSServer         BinaryType = "ssserver"
	BinaryMicrosocks       BinaryType = "microsocks"
	BinarySSHTunUser       BinaryType = "sshtun-user"

	// Client binaries (used in testing)
	BinaryDNSTTClient      BinaryType = "dnstt-client"
	BinarySlipstreamClient BinaryType = "slipstream-client"
	BinarySSLocal          BinaryType = "sslocal"
)

// BinaryDef defines a binary and how to locate it.
type BinaryDef struct {
	Type          BinaryType
	EnvVar        string              // Environment variable for custom path
	PinnedVersion string              // Expected version for this dnstm release
	Platforms     map[string][]string // Supported os -> []arch
}

// DefaultBinaries contains definitions for all supported binaries.
var DefaultBinaries = map[BinaryType]BinaryDef{
	// Server binaries - versions pinned per dnstm release
	BinaryDNSTTServer: {
		Type:   BinaryDNSTTServer,
		EnvVar: "DNSTM_DNSTT_SERVER_PATH",
		Platforms: map[string][]string{
			"linux":   {"amd64", "arm64"},
			"darwin":  {"amd64", "arm64"},
			"windows": {"amd64", "arm64"},
		},
	},
	BinarySlipstreamServer: {
		Type:          BinarySlipstreamServer,
		EnvVar:        "DNSTM_SLIPSTREAM_SERVER_PATH",
		PinnedVersion: "v2026.02.05",
		Platforms: map[string][]string{
			"linux": {"amd64", "arm64"},
		},
	},
	BinarySSServer: {
		Type:          BinarySSServer,
		EnvVar:        "DNSTM_SSSERVER_PATH",
		PinnedVersion: "v1.24.0",
		Platforms: map[string][]string{
			"linux":  {"amd64", "arm64"},
			"darwin": {"amd64", "arm64"},
		},
	},
	BinaryMicrosocks: {
		Type:          BinaryMicrosocks,
		EnvVar:        "DNSTM_MICROSOCKS_PATH",
		PinnedVersion: "v1.0.5",
		Platforms: map[string][]string{
			"linux": {"amd64", "arm64"},
		},
	},
	BinarySSHTunUser: {
		Type:          BinarySSHTunUser,
		EnvVar:        "DNSTM_SSHTUN_USER_PATH",
		PinnedVersion: "v0.3.5",
		Platforms: map[string][]string{
			"linux": {"amd64", "arm64"},
		},
	},

	// Client binaries - pinned versions for testing only
	BinaryDNSTTClient: {
		Type:          BinaryDNSTTClient,
		EnvVar:        "DNSTM_TEST_DNSTT_CLIENT_PATH",
		PinnedVersion: "latest",
		Platforms: map[string][]string{
			"linux":   {"amd64", "arm64"},
			"darwin":  {"amd64", "arm64"},
			"windows": {"amd64", "arm64"},
		},
	},
	BinarySlipstreamClient: {
		Type:          BinarySlipstreamClient,
		EnvVar:        "DNSTM_TEST_SLIPSTREAM_CLIENT_PATH",
		PinnedVersion: "v2026.02.05",
		Platforms: map[string][]string{
			"linux": {"amd64", "arm64"},
		},
	},
	BinarySSLocal: {
		Type:          BinarySSLocal,
		EnvVar:        "DNSTM_TEST_SSLOCAL_PATH",
		PinnedVersion: "v1.23.0",
		Platforms: map[string][]string{
			"linux":  {"amd64", "arm64"},
			"darwin": {"amd64", "arm64"},
		},
	},
}

const (
	// DefaultInstallDir is the default directory for production binaries.
	DefaultInstallDir = "/usr/local/bin"
	// DefaultTestBinDir is the default directory for test binaries.
	DefaultTestBinDir = "tests/.testbin"
)

// Manager handles binary resolution.
type Manager struct {
	binDir string
	os     string
	arch   string
}

// NewManager creates a new binary manager with a specific directory.
func NewManager(binDir string) *Manager {
	return &Manager{
		binDir: binDir,
		os:     runtime.GOOS,
		arch:   runtime.GOARCH,
	}
}

// NewDefaultManager creates a binary manager that auto-detects the environment.
// In test mode, uses tests/.testbin. In production, uses /usr/local/bin.
func NewDefaultManager() *Manager {
	if isTestEnvironment() {
		return NewManager(getTestBinDir())
	}
	return NewManager(DefaultInstallDir)
}

// isTestEnvironment detects if we're running in a test environment.
func isTestEnvironment() bool {
	// Check if running under go test
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}
	// Check if binary name ends with .test
	if strings.HasSuffix(os.Args[0], ".test") {
		return true
	}
	return false
}

// getTestBinDir finds the test binary directory by looking for go.mod.
func getTestBinDir() string {
	dir, _ := os.Getwd()
	for dir != "/" {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, DefaultTestBinDir)
		}
		dir = filepath.Dir(dir)
	}
	return DefaultTestBinDir
}

// GetPath returns the path to an existing binary.
// Resolution order:
// 1. Environment variable (if set and file exists)
// 2. Already in binDir
// Returns error if binary is not found.
func (m *Manager) GetPath(binType BinaryType) (string, error) {
	def, ok := DefaultBinaries[binType]
	if !ok {
		return "", fmt.Errorf("unknown binary type: %s", binType)
	}

	// Check if platform is supported
	if !m.isPlatformSupported(def) {
		return "", fmt.Errorf("binary %s not supported on %s/%s", binType, m.os, m.arch)
	}

	// Check environment variable first
	if def.EnvVar != "" {
		if envPath := os.Getenv(def.EnvVar); envPath != "" {
			if _, err := os.Stat(envPath); err == nil {
				return envPath, nil
			}
			return "", fmt.Errorf("env var %s set to %s but file not found", def.EnvVar, envPath)
		}
	}

	// Check if already in binDir
	binPath := filepath.Join(m.binDir, string(binType))
	if m.os == "windows" {
		binPath += ".exe"
	}
	if _, err := os.Stat(binPath); err == nil {
		return binPath, nil
	}

	return "", fmt.Errorf("binary %s not found at %s (set %s or copy binary to that path)", binType, binPath, def.EnvVar)
}

// EnsureInstalled checks that a binary is available. Returns its path or an error.
func (m *Manager) EnsureInstalled(binType BinaryType) (string, error) {
	return m.GetPath(binType)
}

// EnsureDir creates the binary directory if it doesn't exist.
func (m *Manager) EnsureDir() error {
	return os.MkdirAll(m.binDir, 0755)
}

// BinDir returns the binary directory path.
func (m *Manager) BinDir() string {
	return m.binDir
}

// isPlatformSupported checks if the binary is available for current platform.
func (m *Manager) isPlatformSupported(def BinaryDef) bool {
	archs, ok := def.Platforms[m.os]
	if !ok {
		return false
	}
	for _, a := range archs {
		if a == m.arch {
			return true
		}
	}
	return false
}

// GetDef returns the binary definition for a binary type.
func GetDef(binType BinaryType) (BinaryDef, bool) {
	def, ok := DefaultBinaries[binType]
	return def, ok
}

// CopyToDir copies a binary from srcPath to the manager's binDir.
func (m *Manager) CopyToDir(srcPath string, binType BinaryType) (string, error) {
	if err := m.EnsureDir(); err != nil {
		return "", err
	}

	destName := string(binType)
	if m.os == "windows" {
		destName += ".exe"
	}
	destPath := filepath.Join(m.binDir, destName)

	src, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return destPath, nil
}
