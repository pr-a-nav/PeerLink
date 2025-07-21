package main

import (
	"fmt"
	"log"
	"os"

	// "PeerLink/peer"
	// "PeerLink/torrentfile"
	// "PeerLink/tracker"
)
import "github.com/pr-a-nav/Peerlink/torrentfile"
import "github.com/pr-a-nav/Peerlink/peer"
import "github.com/pr-a-nav/Peerlink/tracker"

func smain() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: peerlink <torrent-file-or-magnet-link>")
		os.Exit(1)
	}
	source := os.Args[1]

	// Parse torrent file or magnet link
	tf, err := torrentfile.New(source)
	if err != nil {
		log.Fatalf("Failed to parse source: %v", err)
	}

	// If magnet, fetch metadata from peers
	if len(tf.PieceHashes) == 0 {
		fmt.Println("Fetching metadata from peers (magnet link)...")
		// Get peers from tracker
		if len(tf.TrackerURLs) == 0 {
			log.Fatal("No trackers found in magnet link")
		}
		peers, err := tracker.GetPeers(tf.TrackerURLs[0], tf.InfoHash, [20]byte{'-', 'P', 'L', '0', '0', '1', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-'}, 6881)
		if err != nil {
			log.Fatalf("Failed to get peers: %v", err)
		}
		if len(peers) == 0 {
			log.Fatal("No peers found from tracker")
		}
		// Connect to first peer and fetch metadata
		p, err := peer.NewClient(peers[0], tf.InfoHash, [20]byte{'-', 'P', 'L', '0', '0', '1', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-'})
		if err != nil {
			log.Fatalf("Failed to connect to peer: %v", err)
		}
		metadata, err := p.FetchMetadata()
		if err != nil {
			log.Fatalf("Failed to fetch metadata: %v", err)
		}
		err = tf.AppendMetadata(metadata)
		if err != nil {
			log.Fatalf("Failed to parse fetched metadata: %v", err)
		}
		fmt.Println("Metadata fetched and parsed.")
	}

	fmt.Printf("Torrent: %s\n", tf.DisplayName)
	fmt.Printf("Trackers: %v\n", tf.TrackerURLs)
	fmt.Printf("Total Length: %d bytes\n", tf.TotalLength)
	fmt.Printf("Pieces: %d\n", len(tf.PieceHashes))

	// Get peers from tracker
	peersList, err := tracker.GetPeers(tf.TrackerURLs[0], tf.InfoHash, [20]byte{'-', 'P', 'L', '0', '0', '1', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-'}, 6881)
	if err != nil {
		log.Fatalf("Failed to get peers: %v", err)
	}
	if len(peersList) == 0 {
		log.Fatal("No peers found from tracker")
	}
	fmt.Printf("Found %d peers.\n", len(peersList))

	// Connect to peers
	var clients []*peer.Client
	for i, addr := range peersList {
		if i >= 10 { // limit to 10 peers for demo
			break
		}
		cli, err := peer.NewClient(addr, tf.InfoHash, [20]byte{'-', 'P', 'L', '0', '0', '1', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-', '-'})
		if err != nil {
			fmt.Printf("Failed to connect to peer %v: %v\n", addr, err)
			continue
		}
		clients = append(clients, cli)
	}
	if len(clients) == 0 {
		log.Fatal("Could not connect to any peers.")
	}
	fmt.Printf("Connected to %d peers.\n", len(clients))

	// Start swarm download
	swarm := peer.NewSwarm(clients, len(tf.PieceHashes), tf.PieceLength, tf.PieceHashes)
	swarm.StartDownload()

	// Collect pieces and write to file
	file, err := os.Create("output.data")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()
	pieces := make([][]byte, len(tf.PieceHashes))
	for i := 0; i < len(tf.PieceHashes); i++ {
		res := <-swarm.DownloadChan
		if res.Err != nil {
			log.Fatalf("Failed to download piece %d: %v", res.Index, res.Err)
		}
		pieces[res.Index] = res.Data
		fmt.Printf("Downloaded piece %d/%d\n", res.Index+1, len(tf.PieceHashes))
	}
	for _, piece := range pieces {
		file.Write(piece)
	}
	fmt.Println("Download complete. Saved to output.data")
}
