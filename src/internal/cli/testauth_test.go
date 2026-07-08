package cli

import (
	"path/filepath"
	"testing"

	"github.com/artschekoff/joplin-cli/src/internal/auth"
)

// withFakeAuth swaps auth's credentials path + device id to test-local values.
func withFakeAuth(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	prevPath, prevID := auth.CredentialsPath, auth.DeviceID
	auth.CredentialsPath = func() (string, error) { return filepath.Join(dir, "credentials"), nil }
	auth.DeviceID = func() (string, error) { return "test-device-id-0123456789", nil }
	t.Cleanup(func() {
		auth.CredentialsPath = prevPath
		auth.DeviceID = prevID
	})
}
