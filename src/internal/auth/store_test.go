package auth

import (
	"errors"
	"path/filepath"
	"testing"
)

func withFakeStore(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	prevPath, prevID := CredentialsPath, DeviceID
	CredentialsPath = func() (string, error) { return filepath.Join(dir, "credentials"), nil }
	DeviceID = func() (string, error) { return "test-device-id-0123456789", nil }
	t.Cleanup(func() { CredentialsPath = prevPath; DeviceID = prevID })
}

func TestSaveLoadRoundTrip(t *testing.T) {
	withFakeStore(t)
	if err := SaveToken("tok-abcdef0123456789"); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := LoadToken()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got != "tok-abcdef0123456789" {
		t.Fatalf("round trip mismatch: %q", got)
	}
}

func TestLoad_NoFile_ReturnsErrNoCredentials(t *testing.T) {
	withFakeStore(t)
	_, err := LoadToken()
	if !errors.Is(err, ErrNoCredentials) {
		t.Fatalf("want ErrNoCredentials, got %v", err)
	}
}

func TestDelete_Idempotent(t *testing.T) {
	withFakeStore(t)
	if err := DeleteToken(); err != nil {
		t.Fatalf("delete on empty: %v", err)
	}
	_ = SaveToken("tok-abcdef0123456789")
	if err := DeleteToken(); err != nil {
		t.Fatalf("delete existing: %v", err)
	}
	if _, err := LoadToken(); !errors.Is(err, ErrNoCredentials) {
		t.Fatalf("expected gone after delete, got %v", err)
	}
}
