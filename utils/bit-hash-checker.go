package utils

import (
	"bytes"
	"crypto/sha1"
)


func BitHashChecker(buf []byte, pieceHash [20]byte) bool {
	hash := sha1.Sum(buf)
	return !bytes.Equal(hash[:], pieceHash[:])
}