package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type clientdb struct {
	group_name []string
	IP         []string
	port       []int
	lastactive time.Time
}

type Message struct {
	ID         string
	sender     string
	content    string
	timestamp  time.Time
	IsReceived bool
}

var messages = []Message{}

func fetchMessages(user string, conn net.Conn) []Message {
	//main work lies here to decode this
	//fetched messages from server for that particular conn
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message:", err)
			return []Message{}
		}
		fmt.Println("Message from server:", msg)
	}
}

func receive(user string, cl clientdb, conn net.Conn) {
	new := fetchMessages(user, conn)
	for _, msg := range new {
		if !msg.IsReceived {
			fmt.Println("New message from:", msg.sender, "Content:", msg.content)
		}
		// a very little logic here is just we double check to make sure we get message after we
		// got inactive
		if msg.timestamp.After(cl.lastactive) {
			newmsg := NewMessages(msg.ID, msg.sender, msg.content, msg.timestamp, true)
			messages = append(messages, newmsg)

		}
	}
}

func NewMessages(ID string, sender string, content string, timestamp time.Time, IsReceived bool) Message {
	return Message{
		ID:         ID,
		sender:     sender,
		content:    content,
		timestamp:  time.Now().UTC(),
		IsReceived: true,
	}
}

func New(ip string, port int, group_name string) clientdb {
	return clientdb{
		group_name: []string{group_name},
		IP:         []string{ip},
		port:       []int{port},
	}
}

func (db clientdb) addgroupname(name string) {
	gname := append(db.group_name, name)
	return
}

func (db clientdb) addIP(ip string) {
	ips := append(db.IP, ip)
	return
}

func (db clientdb) addport(port int) {
	ports := append(db.port, port)
	return
}


func handlemessages(conn  net.Conn){
    defer conn.Close()

    reader := bufio.NewReader(conn)
	message, err := reader.ReadString(byte(reader.Buffered()))
    var msg = append(messages, NewMessages(message.ID, message.sender , message.content , message.timestamp , true))
        
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(message.content)
}



func main() {

	var seradd string = "localhost:9000"
	conn, err := net.Dial("tcp", seradd)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	println("requested")

	// res , err := conn.Write([]byte(" wanna join your server"))
	// if err != nil {
	//     fmt.Println(err)
	// }
	reader := bufio.NewReader(conn)
	group, err := reader.ReadString(byte(reader.Buffered()))
	if err != nil {
		fmt.Println(err)
	}
	println("read group")

	fmt.Println(group)
	fmt.Println("why though")

	fmt.Print("Enter the group you want to join: ")
	scannner := bufio.NewScanner(os.Stdin)
	// fmt.Println(res)
	if scannner.Scan() {

		group = scannner.Text()
		conn.Write([]byte(group))
	}
	confirmation, _ := reader.ReadString('\n')
	fmt.Print(confirmation)

	// Now we can send and receive messages
	for {
		fmt.Print("Enter message or type 'exit' to leave :")
		if scannner.Scan() {
			message := scannner.Text()
			if strings.ToLower(message) == "exit" {
				break
			}
			conn.Write([]byte(message + "\n"))

		}

	}

	conn.Close()
}

func list(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		group, err := reader.ReadString(byte(reader.Buffered()))
		if err != nil {
			break
		}
		fmt.Println(group)
	}

}

func cl_server(){
    listener,err:=net.Listen("tcp","localhost:8000")
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
		go handlemessages(conn,)
	}
}
