package main

import (
	"bufio"
	"fmt"
	"net"
)

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
       
        
    fmt.Println(group)
    // fmt.Println(res)
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