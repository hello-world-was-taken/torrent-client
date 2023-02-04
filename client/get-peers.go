package client

// import "torrent-dsp/model"

// import (
// 	"net"

// )

// func GetPeersFromTracker(torrent *model.Torrent, conn net.Conn) ([]*utils.Peer, error) {
// 	// create a request to send to the tracker
// 	request, err := CreateTrackerRequest(torrent)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// send the request to the tracker
// 	err = SendTrackerRequest(request, conn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// get the response from the tracker
// 	response, err := GetTrackerResponse(conn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// parse the response and get a list of peers
// 	peers, err := ParseTrackerResponse(response)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return peers, nil
// }

// func CreateTrackerRequest(torrent *model.Torrent) (*model.TrackerRequest, error) {
// 	// create a new tracker request
// 	request := model.TrackerRequest{
// 		Info_hash:  torrent.Info_hash,
// 		Peer_id:    "-DS0001-",
// 		Port:       6889,
// 		Uploaded:   0,
// 		Downloaded: 0,
// 		Left:       torrent.Info.Length,
// 		Event:      "started",
// 	}

// 	return &request, nil
// }


