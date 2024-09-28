package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

var sh *shell.Shell

type Entry struct {
	GroupId string
	Port    string
}

func writeToOrbitDB(groupId, port string) {
	entry := Entry{GroupId: groupId, Port: port}
	data := fmt.Sprintf("GroupId: %s, Port: %s", entry.GroupId, entry.Port)

	cid, err := sh.Add(strings.NewReader(data))
	if err != nil {
		log.Fatalf("Failed to write data to IPFS: %v", err)
	}

	fmt.Printf("Data written to IPFS with CID: %s\n", cid)
}

func readFromOrbitDB(cid string) {

	readData, err := sh.Cat(cid)
	if err != nil {
		log.Fatalf("Failed to read data from IPFS: %v", err)
	}

	fmt.Println("Data from IPFS:", readData)
}

func main() {

	sh = shell.NewShell("localhost:5001")

	readCmd := flag.NewFlagSet("read", flag.ExitOnError)
	writeCmd := flag.NewFlagSet("write", flag.ExitOnError)

	writeGroupId := writeCmd.String("groupId", "", "Group ID to write to OrbitDB")
	writePort := writeCmd.String("port", "", "Port to associate with the Group ID")

	readCID := readCmd.String("cid", "", "CID to read data from OrbitDB")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'read' or 'write' subcommands.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "write":
		writeCmd.Parse(os.Args[2:])
		if *writeGroupId == "" || *writePort == "" {
			fmt.Println("Please provide both GroupId and Port.")
			os.Exit(1)
		}
		writeToOrbitDB(*writeGroupId, *writePort)

	case "read":
		readCmd.Parse(os.Args[2:])
		if *readCID == "" {
			fmt.Println("Please provide a CID to read data.")
			os.Exit(1)
		}
		readFromOrbitDB(*readCID)

	default:
		fmt.Println("Unknown command.")
		os.Exit(1)
	}
}
