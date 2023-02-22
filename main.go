package main

import (
	"torrent-dsp/leech"
)

func main() {
	// TODO: remove this hardcoded torrent file names
	// ubuntu-22.10-desktop-amd64.iso
	// vlc-media-player
	// 20A4F6FB1C21B5F5D76BAFDA3D64492125F7FAE2
	// FD6E802C6F3EB1C70367487A55CE3FE782CBC6BC
	// debian-11.6.0-amd64-netinst.iso
	// A00F0AF4A54A170C3A1EE98FA83773E6D73941A5
	// WallpaperDog-10839278.png
	leech.Leech("./torrent-files/debian-11.6.0-amd64-netinst.iso.torrent")
	// infinite loop to keep the program running
	// for { }
}

func prepareSeed(filename string) {
	// open torrent file from the current directory and parse it
}
