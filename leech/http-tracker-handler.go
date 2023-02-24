package leech

import (
	// "fmt"
	"net"
	"net/http"
	"net/url"
	"torrent-dsp/constant"
	"torrent-dsp/model"
	"torrent-dsp/utils"

	"github.com/zeebo/bencode"
)

// Get peers from one of the trackers in the torrent file
func GetPeersFromTrackers(torrent *model.Torrent) ([]model.Peer, error) {

	// create a request to send to the tracker
	httpTrackerURLs, err := buildTrackerRequestURLs(torrent)
	if err != nil {
		return nil, err
	}

	peers, err := getPeersFromTrackersHelper(httpTrackerURLs)
	if err != nil {
		return nil, err
	}

	peers = []model.Peer{model.Peer{IP: net.IP([]byte{192, 168, 43, 5}), Port: 6881}}

	return peers, nil
}

// builds a tracker request url from the announcement list
func buildTrackerRequestURLs(torrent *model.Torrent) ([]string, error) {

	// query parameters to the tracker url
	requestParams := model.TrackerRequestParams{
		Info_hash:  torrent.InfoHash,
		Peer_id:    constant.CLIENT_ID,
		Port:       6881,
		Uploaded:   0,
		Downloaded: 0,
		Left:       torrent.Info.Length,
		Compact:    1,
		Event:      "started",
	}

	// list of urls
	URLs := []string{}

	// check if the announce is an http tracker
	if utils.IsHTTPTracker(torrent.Announce) {
		URL, err := url.Parse(torrent.Announce)
		if err != nil {
			return []string{}, err
		}
		URL.RawQuery = requestParams.Encode()
		URLs = append(URLs, URL.String())
	}

	// go through all the trackers in the announce list and create a request for each one
	// discard udp and only take http trackers
	for _, tracker := range torrent.AnnounceList {

		// parse the tracker url
		URL, err := url.Parse(tracker[0])
		if err != nil {
			return []string{}, err
		}

		if utils.IsHTTPTracker(URL.String()) {
			// add the query parameters to the tracker url
			URL.RawQuery = requestParams.Encode()
			URLs = append(URLs, URL.String())
		}
	}

	return URLs, nil
}

// go through each HTTP tracker requests until we get a response from one
func getPeersFromTrackersHelper(URLs []string) ([]model.Peer, error) {
	peers := []model.Peer{}
	for _, URL := range URLs {
		response, err := getPeerFromURL(URL)
		if err == nil && len(response) > 0 {
			for _, p := range response {
				peers = append(peers, p)
			}
		}
	}

	return peers, nil
}

// try to get peers from a passed in tracker request url
func getPeerFromURL(URL string) ([]model.Peer, error) {

	// send the request to the tracker
	response, err := http.Get(URL)
	if err != nil {
		return []model.Peer{}, err
	}

	// close the response body
	defer response.Body.Close()

	// decode the response
	trackerResponse := model.TrackerResponse{}
	err = bencode.NewDecoder(response.Body).Decode(&trackerResponse)
	if err != nil {
		return []model.Peer{}, err
	}

	peers, err := model.PeerParser([]byte(trackerResponse.Peers))
	if err != nil {
		return []model.Peer{}, err
	}

	// log the peers
	// utils.LogPeers(peers)

	// fmt.Println("Retrieved peers successfully")
	return peers, nil
}
