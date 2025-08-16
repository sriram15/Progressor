package profile_test

import (
	"testing"

	"github.com/sriram15/progressor-todo-app/internal/profile"
)

func TestKeyringStorage(t *testing.T) {
	const testKey = "test_profile_for_go_test"
	const testToken = "a-secret-token-from-go-test"

	// Defer cleanup to ensure the token is deleted even if the test fails.
	t.Cleanup(func() {
		t.Log("Cleaning up by deleting the token...")
		err := profile.DeleteToken(testKey)
		// It is not an error if the secret is not found during cleanup,
		// as it may have been deleted by the test itself.
		if err != nil && err.Error() != "secret not found in keyring" {
			t.Errorf("Cleanup failed with an unexpected error: %v", err)
		}
	})

	// 1. Store the token
	t.Log("Attempting to store the token...")
	err := profile.StoreToken(testKey, testToken)
	if err != nil {
		t.Fatalf("Failed to store token: %v", err)
	}
	t.Log("Successfully stored the token.")

	// 2. Retrieve the token
	t.Log("Attempting to retrieve the token...")
	retrievedToken, err := profile.GetToken(testKey)
	if err != nil {
		t.Fatalf("Failed to retrieve token: %v", err)
	}

	// 3. Verify the token
	if retrievedToken != testToken {
		t.Fatalf("Verification failed: Retrieved token \"%s\" does not match original token \"%s\"", retrievedToken, testToken)
	}
	t.Logf("Successfully retrieved and verified token: %s", retrievedToken)

	// 4. Explicitly delete for verification step (cleanup will also run)
	t.Log("Deleting token to verify deletion...")
	err = profile.DeleteToken(testKey)
	if err != nil {
		t.Fatalf("Failed to delete token for verification: %v", err)
	}

	// 5. Verify deletion
	t.Log("Verifying deletion...")
	_, err = profile.GetToken(testKey)
	if err == nil {
		t.Fatalf("Verification failed: Token should have been deleted but was found.")
	}
	t.Log("Successfully verified that the token is deleted.")
}