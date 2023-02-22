package utils

import (
	"fmt"
	"bytes"
	"crypto/sha1"
)


func BitHashChecker(buf []byte, pieceHash [20]byte) bool {
	hash := sha1.Sum(buf)
	if bytes.Equal(hash[:], pieceHash[:]) {
		return true
	}
	fmt.Println("Hash mismatch")
	return false
}