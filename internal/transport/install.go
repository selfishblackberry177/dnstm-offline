package transport

import (
	"fmt"

	"github.com/net2share/dnstm/internal/binary"
	"github.com/net2share/dnstm/internal/config"
	"github.com/net2share/dnstm/internal/log"
)

// StatusFunc is a callback for reporting installation status messages.
type StatusFunc func(message string)

// EnsureTransportBinariesInstalled checks that required binaries for a transport type are present.
func EnsureTransportBinariesInstalled(transport config.TransportType) error {
	switch transport {
	case config.TransportSlipstream:
		return EnsureSlipstreamInstalled()
	case config.TransportDNSTT:
		return EnsureDnsttInstalled()
	default:
		return nil
	}
}

// EnsureBackendBinariesInstalled checks that required binaries for a backend type are present.
func EnsureBackendBinariesInstalled(backend config.BackendType) error {
	switch backend {
	case config.BackendShadowsocks:
		return EnsureShadowsocksInstalled()
	default:
		return nil
	}
}

// EnsureDnsttInstalled checks that dnstt-server is present.
func EnsureDnsttInstalled() error {
	return EnsureDnsttInstalledWithStatus(nil)
}

// EnsureDnsttInstalledWithStatus checks that dnstt-server is present with status callback.
func EnsureDnsttInstalledWithStatus(statusFn StatusFunc) error {
	return ensureBinaryPresent(binary.BinaryDNSTTServer, "dnstt-server", statusFn)
}

// EnsureSlipstreamInstalled checks that slipstream-server is present.
func EnsureSlipstreamInstalled() error {
	return EnsureSlipstreamInstalledWithStatus(nil)
}

// EnsureSlipstreamInstalledWithStatus checks that slipstream-server is present with status callback.
func EnsureSlipstreamInstalledWithStatus(statusFn StatusFunc) error {
	return ensureBinaryPresent(binary.BinarySlipstreamServer, "slipstream-server", statusFn)
}

// EnsureShadowsocksInstalled checks that ssserver is present.
func EnsureShadowsocksInstalled() error {
	return EnsureShadowsocksInstalledWithStatus(nil)
}

// EnsureShadowsocksInstalledWithStatus checks that ssserver is present with status callback.
func EnsureShadowsocksInstalledWithStatus(statusFn StatusFunc) error {
	return ensureBinaryPresent(binary.BinarySSServer, "ssserver", statusFn)
}

// EnsureSSHTunUserInstalled checks that sshtun-user is present.
func EnsureSSHTunUserInstalled() error {
	return EnsureSSHTunUserInstalledWithStatus(nil)
}

// EnsureSSHTunUserInstalledWithStatus checks that sshtun-user is present with status callback.
func EnsureSSHTunUserInstalledWithStatus(statusFn StatusFunc) error {
	return ensureBinaryPresent(binary.BinarySSHTunUser, "sshtun-user", statusFn)
}

// IsSSHTunUserInstalled checks if sshtun-user binary is installed.
func IsSSHTunUserInstalled() bool {
	mgr := binary.NewDefaultManager()
	_, err := mgr.GetPath(binary.BinarySSHTunUser)
	return err == nil
}

// ensureBinaryPresent checks that a binary exists using the binary manager.
func ensureBinaryPresent(binType binary.BinaryType, displayName string, statusFn StatusFunc) error {
	mgr := binary.NewDefaultManager()

	path, err := mgr.GetPath(binType)
	if err != nil {
		return fmt.Errorf("%s not found: %w", displayName, err)
	}

	log.Debug("%s found at %s", displayName, path)

	if statusFn != nil {
		statusFn(fmt.Sprintf("%s ready", displayName))
	}
	return nil
}
