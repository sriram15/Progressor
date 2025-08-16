package profile

import "github.com/zalando/go-keyring"

// serviceName is the identifier for our application's secrets in the OS keychain.
const serviceName = "com.progressor.profiles"

// StoreToken saves a token securely in the OS keychain.
// The key is typically the profile ID.
func StoreToken(key, token string) error {
	return keyring.Set(serviceName, key, token)
}

// GetToken retrieves a token from the OS keychain.
// The key is typically the profile ID.
func GetToken(key string) (string, error) {
	return keyring.Get(serviceName, key)
}

// DeleteToken removes a token from the OS keychain.
// The key is typically the profile ID.
func DeleteToken(key string) error {
	return keyring.Delete(serviceName, key)
}
