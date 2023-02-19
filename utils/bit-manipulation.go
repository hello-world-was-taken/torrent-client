package utils

func BitOn(byteField []byte, bitIdx int) bool {
	// first get the byte index
	byteIdx := bitIdx / 8
	bitIdxOnCurrByte := bitIdx % 8

	// on bit means it has the piece data
	// check if the length of the bit field is greater than the byte index
	// if not, then the peer does not have the piece
	if len(byteField) <= byteIdx {
		return false
	}
	return byteField[byteIdx] & (1 << bitIdxOnCurrByte) != 0
}


func TurnBitOn(byteField []byte, bitIdx int) {
	// first get the byte index
	byteIdx := bitIdx / 8
	bitIdxOnCurrByte := bitIdx % 8

	byteField[byteIdx] |= 1 << bitIdxOnCurrByte
}