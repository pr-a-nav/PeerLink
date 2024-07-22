package main

import ("fmt"
        "net"
		// "net/http"
		

		
)

func server(){
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
	}

}

func handle(conn net.Conn){
  defer conn.Close()

   add := conn.RemoteAddr().(*net.TCPAddr)
   clientIP := add.IP.String()
   clientport := add.Port

   fmt.Printf("IP %s port %d", clientIP,clientport)
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
   Members map[string]*Member
   messages []Message
}

type Message struct{
	SenderID string
	Content string
}

type Member struct{
    ID string
	Port int
}


func NewIpad(ip string, port int) Ipad {
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
	res := NewIpad("12.45.56", 45)
	res.add("12.4.6", 45)
}