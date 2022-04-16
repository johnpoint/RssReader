package utils

import (
	"crypto/sha256"
	"fmt"
)

func Sha256(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
