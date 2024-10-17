package main

import (
	"bufio"
	"fmt"
	"net"
	 orbitdb "github.com/orbitdb/go-orbit-db"
	"github.com/ipfs/go-ipfs/core"
	// "bufio"
	// "context"
	// "encoding/json"
	"sync"
	// "github.com/pr-a-nav/Peerlink/orbitdb"
	// "net/http"
)

var dbLock sync.Mutex

type peerID = map[string]string

type Peer struct {
	Address string
	Port    int
}

var groups = map[string][]Peer{
	"SDE 1":        {},
	"Ml Guys":      {},
	"Torrent Guys": {},
}

type RendezvousPoint struct {
	identifier int
	peerlist   []peerID
	mu         sync.Mutex
}

func CollectDataFromIP(rp *RendezvousPoint, ipAddr string, data interface{}) error {

	// rp.peerlist[ipAddr] = data

	// Store updated rendezvous point in OrbitDB
	// err := rm.db.Put(rp.Identifier, rp)
	// if err != nil {
	// 		return fmt.Errorf("failed to update rendezvous point in OrbitDB: %w", err)
	// }

	return nil
}
func server() {
	fmt.Println("server started")
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("error", err)
	}
	for {
		req, err := ln.Accept()
		if err != nil {
			fmt.Println("error", err)
		}
		// reader  := bufio.NewReader(req)
		// message ,err := reader.ReadString(byte(reader.Buffered()))
		// fmt.Println(message)
		list := "groups"
		for group := range groups {
			list += group + "\n"
		}
		req.Write([]byte(list))

		handle(req)
	}

}

func handle(conn net.Conn) {
	defer conn.Close()
	res := New("12.45.56", 45)

	add := conn.RemoteAddr().(*net.TCPAddr)
	clientIP := add.IP.String()
	clientport := add.Port
	register(res, clientIP, clientport)
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString(byte(reader.Buffered()))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(message)
	//    var selected string
	//    _ , err := fmt.Fscanln(conn, &selected)
	//    if err!=nil{
	// 	fmt.Print(err)
	//    }
	//    fmt.Println(selected)
	fmt.Println("nopes")
	fmt.Printf("IP %s port %d", clientIP, clientport)
	//    list := "groups"
	// 	for  group := range groups{
	// 		list +=group+"\n"
	// 	}
	// 	conn.Write([]byte(list))

}

func sendgroups(conn net.Conn) {
	list := "groups"
	for group := range groups {
		list += group + "\n"
	}
	conn.Write([]byte(list))
}

func register(Ipadd Ipad, ip string, port int) {
	Ipad.add(Ipadd, ip, port)

}

func vishwas() {
	fmt.Println("hello pranav")
}

type Ipad struct {
	ipad []string
	port []int
}

type Group struct {
	ID       string
	Name     string
	Members  map[string]*Ipad
	messages []Message
}

type Message struct {
	SenderID string
	Content  string
}

func New(ip string, port int) Ipad {
	return Ipad{
		ipad: []string{ip},
		port: []int{port},
	}
}
func (i Ipad) add(ipad string, port int) {
	i.ipad = append(i.ipad, ipad)
	i.port = append(i.port, port)

}
func setupOribitDB()(*core.IpfsNode, orbitdb.OrbitDB){
	node,err:=ipfs.NewIPFSNode()
	if err!=nil{
		fmt.Print("Error in IPFS :",err)
	}

	db, err:=orbitdb.NewOrbitDB(node)

	if err!=nil{
		fmt.Print("Error initialising orbitdb:",err)
	}
	return node,db
}
func main() {
	res := New("12.45.56", 45)
	res.add("12.4.6", 45)
	node, db:=setupOribitDB()

	listener,err:=net.Listen("tcp","localhost:9000")
	if err!=nil{
		fmt.Println("Error starting",err)
		return
	}

	defer listener.Close()

	fmt.Println(("Server is listening on host:9000"))

	for{
		conn,err:=listener.Accept()
		if err!=nil{
			fmt.Println(("error accepting connection:",err))
		}
		go handlemessages(conn,db)
	}

	// server() included this fucntion in main itself
}


