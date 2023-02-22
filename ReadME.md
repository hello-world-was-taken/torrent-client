# Golang Torrent Client
This is a command-line torrent client application developed in Golang that can download and upload torrent files. The application allows users to download files from the BitTorrent network and seed files to share them with other users. It uses the standard BitTorrent protocol to connect to other peers and download or upload files.
# Building from Source
Clone the repository using the following command:

```sh
git clone https://github.com/your-username/your-repo.git
```

Navigate to the directory containing the source code:
```sh
cd your-repo
```

Run the application using the following command:

```sh
go run main.go
```
The client will start downloading the file as soon as it is launched. You need to first download the file before attempting to seed.

# Functionality
The Golang Torrent Client has the following features:

- ***Torrent Parser***: Parses the metadata of a torrent file and extracts information such as the file name, file size, and list of trackers.
- ***Seeder***: Seeds a torrent by connecting to other peers and sharing the file with them.
- ***Leecher***: Downloads a torrent by connecting to other peers and downloading the file from them.
- ***Client Factory***: Creates a seeder or leecher client depending on the mode specified by the user.

# Credits
This project was developed by:
- ***Hailemariam Arega***
- ***Nazrawi Demeke***
- ***Biruk Tassew***
- ***Biruk Anley***
- ***Tamirat Dereje***