package main

import (
	"time"
	"torrent-dsp/leech"
	"torrent-dsp/seed"
)

func main() {
	go seed.SeederMain()
	time.Sleep(5 * time.Second)
	leech.Leech("./torrent-files/debian-11.6.0-amd64-netinst.iso.torrent")
}
