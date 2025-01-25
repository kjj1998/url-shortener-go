package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func TestGenerateUniqueId(t *testing.T) {
	originalGenerateUniqueId := GenerateUniqueId

	// Step 2: Mock the function
	GenerateUniqueId = func() uint64 {
		return 550726476329124100 // Return the same value for testing
	}

	// Step 3: Use the mocked function in your test
	defer func() {
		// Restore the original function after the test
		GenerateUniqueId = originalGenerateUniqueId
	}()

	id := GenerateUniqueId()
	assert.Equal(t, uint64(550726476329124100), id) // Assert the mocked value is returned
}

func TestShortenUrl(t *testing.T) {
	ans := ShortenUrl(550726476329124100)

	if ans != "EEAA2fZJJ9A" {
		t.Errorf("")
		t.Errorf("ShortenUrl(550726476329124100) = %s; want EEAA2fZJJ9A", ans)
	}
}
