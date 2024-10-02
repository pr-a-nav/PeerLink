package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
    "time"
)


type clientdb struct{
    group_name []string 
    IP []string
    port []int
    lastactive time.Time

}

type Message struct {
    ID          string
    sender      string
    content     string
    timestamp   time.Time
    IsReceived  bool 
}

var messages =[]Message{}

func fetchMessages(user string, conn net.Conn) []Message {
    reader := bufio.NewReader(conn)
    group,err :=reader.ReadString(byte(reader.Buffered()))
    return []Message{} 
}

func receive(user string , cl clientdb , conn net.Conn) {
    new := fetchMessages(user , conn)
    for _, msg := range new{
        if !msg.IsReceived {
            fmt.Println("New message from:", msg.sender, "Content:", msg.content)
        }
        if msg.timestamp > cl.lastactive{ 
           newmsg := NewMessages(msg.ID,msg.sender,msg.content,msg.timestamp,true)
           messages = append(messages, newmsg)
          
        }
    }
}

func NewMessages( ID string,sender string, content string,timestamp  time.Time, IsReceived  bool ) Message {
    return Message{
        ID : ID,
        sender :sender,
        content : content,
        timestamp: time.Now().UTC(),
        IsReceived: true,

    }
}

func New(ip string, port int , group_name string) clientdb {
    return clientdb{
        group_name: []string{group_name},
        IP:   []string{ip},  
        port: []int{port}, 
    }
}

func (db clientdb ) addgroupname( name string) {
   gname := append(db.group_name, name)
   return 
}

func (db clientdb ) addIP( ip string) {
    ips := append(db.IP, ip)
    return 
 }

 func (db clientdb ) addport( port int) {
    ports := append(db.port, port)
    return 
 }

func main() {
	
	var seradd string = "localhost:9000"
    conn, err := net.Dial("tcp", seradd)
    if err != nil {
        fmt.Println(err)
    }
    println("requested")

    // res , err := conn.Write([]byte(" wanna join your server"))
    // if err != nil {
    //     fmt.Println(err)
    // }
    reader := bufio.NewReader(conn)
    group,err :=reader.ReadString(byte(reader.Buffered()))
    if err !=nil{
        fmt.Println(err)
    }
    println("read group")

        
    fmt.Println(group)
    fmt.Println("why though")

    scannner := bufio.NewScanner(os.Stdin)
    // fmt.Println(res)
    if scannner.Scan(){

    group =scannner.Text()
    conn.Write([]byte(group))
}

    conn.Close()
}

func list(conn net.Conn){
    reader := bufio.NewReader(conn)
    for {
        group,err :=reader.ReadString(byte(reader.Buffered()))
        if err!= nil{
            break
        }
    fmt.Println(group)}

}