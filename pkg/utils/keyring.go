package utils

import (
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	keyService = "urlscan/urlscan-cli"
	keyName    = "URLSCAN_API_KEY"
)

type KeyManager struct{}

func NewKeyManager() *KeyManager {
	return &KeyManager{}
}

func (tm *KeyManager) GetKey() (string, error) {
	s, err := keyring.Get(keyService, keyName)
	if err != nil {
		return "", fmt.Errorf("get a urlscan.io API key from keyring: %w", err)
	}
	return s, nil
}

func (tm *KeyManager) SetKey(token string) error {
	if err := keyring.Set(keyService, keyName, token); err != nil {
		return fmt.Errorf("set a urlscan.io API key in keyring: %w", err)
	}
	return nil
}

func (tm *KeyManager) RemoveKey() error {
	err := keyring.Delete(keyService, keyName)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			fmt.Println("API key not found in keyring, nothing to delete.")
			return nil
		}
		return fmt.Errorf("delete a urlscan.io API key from keyring: %w", err)
	}
	return nil
}

func (tm *KeyManager) CheckSercretService() error {
	_, err := keyring.Get(keyService, keyName)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return nil // Keyring service is available but no key is set
		}
		return errors.Join(fmt.Errorf("keyring service is unavialable, check your keyring configuration"), err)
	}
	return nil // Keyring service is available and a key exists
}
