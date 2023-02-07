package constant

import "time"

// client ID is a 20 byte string
const (
	CLIENT_ID = "-TR2940-k8hj0wgej6ch"
)

// connection timeout is the time after which the connection is closed
const (
	CONNECTION_TIMEOUT = 3 * time.Second
)


// MessageID is the type of message that is sent/received over the wire
const (
	CHOKE = 0
	UNCHOKE = 1
	INTERESTED = 2
	NOT_INTERESTED = 3
	HAVE = 4
	BITFIELD = 5
	REQUEST = 6
	PIECE = 7
	CANCEL = 8
)
