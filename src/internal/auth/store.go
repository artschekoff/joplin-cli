package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/denisbrodbeck/machineid"
)

const AppName = "joplin-cli"

// ErrNoCredentials is returned when the credentials file is absent.
var ErrNoCredentials = errors.New("no credentials found — run `joplin-cli login` to store your token")

// Overridable in tests. Restore via t.Cleanup.
var (
	CredentialsPath = defaultCredentialsPath
	DeviceID        = defaultDeviceID
)

func defaultCredentialsPath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("user config dir: %w", err)
	}
	return filepath.Join(base, AppName, "credentials"), nil
}

func defaultDeviceID() (string, error) {
	id, err := machineid.ID()
	if err != nil {
		return "", fmt.Errorf("read machine id: %w", err)
	}
	return id, nil
}

// newAEAD builds an AES-256-GCM cipher using SHA-256(deviceID) as the key.
func newAEAD() (cipher.AEAD, error) {
	id, err := DeviceID()
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256([]byte(id))
	block, err := aes.NewCipher(sum[:])
	if err != nil {
		return nil, fmt.Errorf("aes: %w", err)
	}
	return cipher.NewGCM(block)
}

// SaveToken encrypts the token with AES-256-GCM and writes it with mode 0600.
func SaveToken(token string) error {
	if token == "" {
		return errors.New("token is empty")
	}
	gcm, err := newAEAD()
	if err != nil {
		return err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("nonce: %w", err)
	}
	ct := gcm.Seal(nil, nonce, []byte(token), nil)
	out := append(nonce, ct...)

	path, err := CredentialsPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	if err := os.WriteFile(path, out, 0o600); err != nil {
		return fmt.Errorf("write credentials: %w", err)
	}
	return nil
}

// LoadToken reads and decrypts the credentials file.
func LoadToken() (string, error) {
	path, err := CredentialsPath()
	if err != nil {
		return "", err
	}
	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", ErrNoCredentials
	}
	if err != nil {
		return "", fmt.Errorf("read credentials: %w", err)
	}
	gcm, err := newAEAD()
	if err != nil {
		return "", err
	}
	if len(raw) < gcm.NonceSize() {
		return "", fmt.Errorf("credentials file corrupt (too short)")
	}
	nonce, ct := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt credentials: %w (device changed or file tampered — run `joplin-cli login` again)", err)
	}
	return string(pt), nil
}

// DeleteToken removes the credentials file. Idempotent — no error when absent.
func DeleteToken() error {
	path, err := CredentialsPath()
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
