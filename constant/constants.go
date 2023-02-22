package constant

import "time"

// client ID is a 20 byte string
const (
	// generate another client id
	CLIENT_ID = "-TR3000-GvlaOyQjADkn"
)

// connection timeout is the time after which the connection is closed
const (
	CONNECTION_TIMEOUT = 10 * time.Second
	PIECE_DOWNLOAD_TIMEOUT = 30 * time.Second
	PIECE_UPLOAD_TIMEOUT = 30 * time.Second
	MAX_RETRY_COUNT = 5
)

// MessageID is the type of message that is sent/received over the wire
const (
	CHOKE          = 0
	UN_CHOKE       = 1
	INTERESTED     = 2
	NOT_INTERESTED = 3
	HAVE           = 4
	BIT_FIELD      = 5
	REQUEST        = 6
	PIECE          = 7
	CANCEL         = 8
)


// batch block download 
const (
	MAX_BATCH_DOWNLOAD = 5
	MAX_BLOCK_LENGTH = 16384 // 16KB
)


// seeder port
const (
	SEEDER_PORT = 6881
)
