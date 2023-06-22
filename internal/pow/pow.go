package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// HashcashData - struct with fields of Hashcash
type HashcashData struct {
	ZerosCount int
	Resource   string
	Counter    int
}

// ToString - stringifies hashcash for next sending it on TCP
func (h HashcashData) ToString() string {
	builder := strings.Builder{}
	builder.WriteString(strconv.Itoa(h.ZerosCount))
	builder.WriteString(":")
	builder.WriteString(h.Resource)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(h.Counter))
	return builder.String()
}

// sha1Hash - calculates sha1 hash from given string
func sha256Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// IsHashCorrect - checks that hash has leading <zerosCount> zeros
func IsHashCorrect(hash string, zerosCount int) bool {
	if zerosCount > len(hash) {
		return false
	}
	for i := 0; i < zerosCount; i++ {
		if hash[i] != '0' {
			return false
		}
	}

	return true
}

// BruteForceHashcash - calculates correct hashcash by bruteforce
// until the resulting hash satisfies the condition of IsHashCorrect
// maxIterations to prevent endless computing (0 or -1 to disable it)
func (h HashcashData) BruteForceHashcash(maxIterations int) (HashcashData, error) {
	for {
		header := h.ToString()
		hash := sha256Hash(header)
		if IsHashCorrect(hash, h.ZerosCount) {
			return h, nil
		}
		h.Counter++

		if maxIterations > 0 && h.Counter > maxIterations {
			break
		} else if maxIterations <= 0 {
			continue
		}
	}

	return h, fmt.Errorf("max iterations exceeded")
}
