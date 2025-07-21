# Peerlink

Peerlink is a decentralized social network and BitTorrent client. It allows users to:
- Discover peers and join groups via a simple TCP server
- Chat with other users in groups
- Download files using the BitTorrent protocol (from .torrent files or magnet links)

## Features
- **Peer Discovery Server:** Maintains a list of connected peers and groups.
- **Client:** Connects to the server, joins groups, and sends/receives messages.
- **BitTorrent Client:** Parses .torrent files or magnet links, discovers peers, downloads files in pieces, and assembles the final file.

---

## Prerequisites
- Go 1.18 or higher

---

## Installation
1. Clone the repository:
   ```sh
   git clone https://github.com/pr-a-nav/Peerlink.git
   cd Peerlink
   ```
2. Download dependencies:
   ```sh
   go mod tidy
   ```

---

## Usage

### 1. Start the Peer Discovery Server
In one terminal:
```sh
go run server.go
```
You should see:
```
Server is listening on :9000
```

### 2. Start a Client
In another terminal:
```sh
go run client/client.go
```
- The client will connect to the server, print the peer list, and prompt you to join a group and send messages.

### 3. Start the BitTorrent Client
In a third terminal:
```sh
go run peerlink.go <torrent-file-or-magnet-link>
```
- Replace `<torrent-file-or-magnet-link>` with the path to a .torrent file or a magnet link.
- The client will parse the file/link, discover peers, and download the file as `output.data`.

---

## Example Workflow
1. **Alice** runs the server: `go run server.go`
2. **Bob** runs the client: `go run client/client.go`
   - Bob sees the peer list, joins a group, and sends a message.
3. **Alice** runs the BitTorrent client: `go run peerlink.go myfile.torrent`
   - Alice downloads a file using BitTorrent, possibly using the peer list from the server for additional peer discovery.

---

## Project Structure
```
Peerlink/
├── client/         # Client code for group chat and peer discovery
│   └── client.go
├── peer/           # BitTorrent peer protocol implementation
├── torrentfile/    # Torrent file and magnet link parsing
├── tracker/        # Tracker communication
├── server.go       # Peer discovery/group server
├── peerlink.go     # BitTorrent client main file
└── README.md
```

---

## Contributing
Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

---

## License
[MIT](LICENSE) 