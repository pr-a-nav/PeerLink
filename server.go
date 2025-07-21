package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

var dbLock sync.Mutex

type Peer struct {
	Address string
	Port    int
}

var groups = map[string][]Peer{
	"SDE 1":        {},
	"Ml Guys":      {},
	"Torrent Guys": {},
}

type Ipad struct {
	ipad []string
	port []int
}

func New(ip string, port int) Ipad {
	return Ipad{
		ipad: []string{ip},
		port: []int{port},
	}
}

func (i *Ipad) add(ipad string, port int) {
	i.ipad = append(i.ipad, ipad)
	i.port = append(i.port, port)
}

var peerListMu sync.Mutex
var peerList []Peer

func handle(conn net.Conn) {
	defer conn.Close()
	add := conn.RemoteAddr().(*net.TCPAddr)
	clientIP := add.IP.String()
	clientport := add.Port
	res := New(clientIP, clientport)
	register(&res, clientIP, clientport)

	// Add to global peer list
	peerListMu.Lock()
	peerList = append(peerList, Peer{Address: clientIP, Port: clientport})
	// Prepare peer list string
	peerListStr := "Peers:\n"
	for _, p := range peerList {
		peerListStr += fmt.Sprintf("%s:%d\n", p.Address, p.Port)
	}
	peerListMu.Unlock()

	// Send peer list to the client
	conn.Write([]byte(peerListStr))
	conn.Write([]byte("\n"))

	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Received message:", message)
	fmt.Printf("IP %s port %d\n", clientIP, clientport)
}

func register(ipadd *Ipad, ip string, port int) {
	ipadd.add(ip, port)
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on :9000")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handle(conn)
	}
}
