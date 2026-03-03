package binary

import (
	"os"
	"testing"
)

func TestGetPath_EnvVarOverride(t *testing.T) {
	mgr := NewManager(t.TempDir())

	// Create a fake binary
	tmpFile := t.TempDir() + "/fake-binary"
	if err := os.WriteFile(tmpFile, []byte("fake"), 0755); err != nil {
		t.Fatal(err)
	}

	// Set env var
	os.Setenv("DNSTM_TEST_DNSTT_CLIENT_PATH", tmpFile)
	defer os.Unsetenv("DNSTM_TEST_DNSTT_CLIENT_PATH")

	path, err := mgr.GetPath(BinaryDNSTTClient)
	if err != nil {
		t.Fatalf("GetPath failed: %v", err)
	}

	if path != tmpFile {
		t.Errorf("Expected %s, got %s", tmpFile, path)
	}
}

func TestIsPlatformSupported(t *testing.T) {
	mgr := NewManager(t.TempDir())

	// DNSTT should be supported on current platform (linux/darwin/windows)
	def := DefaultBinaries[BinaryDNSTTClient]
	if !mgr.isPlatformSupported(def) {
		t.Errorf("Expected DNSTT to be supported on %s/%s", mgr.os, mgr.arch)
	}
}
