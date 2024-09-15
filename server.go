package main

import ("fmt"
        "net"
		// "context"
        // "encoding/json"
        "sync"
		// "github.com/pr-a-nav/Peerlink/orbitdb"
		// "net/http"
		

		
)


 type peerID  = map[string]string


 type RendezvousPoint struct {
	identifier  int
	peerlist  []peerID
	    mu   sync.Mutex

 }
 func  CollectDataFromIP(rp *RendezvousPoint, ipAddr string, data interface{}) error {
	

	// rp.peerlist[ipAddr] = data

	// Store updated rendezvous point in OrbitDB
	// err := rm.db.Put(rp.Identifier, rp)
	// if err != nil {
	// 		return fmt.Errorf("failed to update rendezvous point in OrbitDB: %w", err)
	// }

	return nil
}
func server(){
	fmt.Println("server started")
	ln ,err := net.Listen("tcp", ":9000")
	 if err!=nil{
		fmt.Println("error", err)
	 }
	for {
		req ,err := ln.Accept()
		if err!=nil{
			fmt.Println("error", err)
		}
		fmt.Println(req)
		go handle(req)
	}

}

func handle(conn net.Conn){
  defer conn.Close()
   res := New("12.45.56", 45)

   add := conn.RemoteAddr().(*net.TCPAddr)
   clientIP := add.IP.String()
   clientport := add.Port
   register(res, clientIP, clientport)

   fmt.Printf("IP %s port %d", clientIP,clientport)
}

func register(Ipadd Ipad, ip string , port int){
	Ipad.add(Ipadd, ip, port)

}


func vishwas() {
	fmt.Println("hello pranav")
}

type Ipad struct{
    ipad []string
	port []int
}

type Group struct{
   ID string
   Name string
   Members map[string]*Ipad
   messages []Message
}

type Message struct{
	SenderID string
	Content string
}



func New(ip string, port int) Ipad {
    return Ipad{
        ipad:   []string{ip},  
        port: []int{port}, 
    }
}
func (i Ipad )add(ipad string, port int ){
i.ipad = append(i.ipad,ipad )
i.port = append(i.port, port)

}

func main(){
	res := New("12.45.56", 45)
	res.add("12.4.6", 45)
	server()
}