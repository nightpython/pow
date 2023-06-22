package pow

import (
	"testing"
)

func TestHashcashData(t *testing.T) {
	t.Run("ToString", func(t *testing.T) {
		hashcash := HashcashData{
			ZerosCount: 3,
			Resource:   "example.com",
			Counter:    1,
		}

		expectedString := "3:example.com:1"
		result := hashcash.ToString()

		if result != expectedString {
			t.Errorf("Expected string: %s, but got: %s", expectedString, result)
		}
	})

	t.Run("IsHashCorrect", func(t *testing.T) {
		correctHash := "000abcd"
		incorrectHash := "00abcd"

		if !IsHashCorrect(correctHash, 3) {
			t.Errorf("Expected correct hash to be valid")
		}

		if IsHashCorrect(incorrectHash, 3) {
			t.Errorf("Expected incorrect hash to be invalid")
		}
	})

	t.Run("BruteForceHashcash", func(t *testing.T) {
		hashcash := HashcashData{
			ZerosCount: 5,
			Resource:   "example.com",
			Counter:    1,
		}

		// Test successful case
		result, err := hashcash.BruteForceHashcash(1000000)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
		if !IsHashCorrect(sha256Hash(result.ToString()), 5) {
			t.Errorf("Expected computed hashcash to have correct hash")
		}

		// Test case where max iterations are exceeded
		_, err = hashcash.BruteForceHashcash(1)
		if err == nil {
			t.Errorf("Expected error due to max iterations exceeded")
		}
	})

}
