package utils

func ConvertStringToByteArray(str string) *[20]byte {
	var bytes [20]byte
	copy(bytes[:], []byte(str))
	return &bytes
}