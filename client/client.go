package main

import (
    "fmt"
    "net"
)

func main() {
	
	var seradd string = "localhost:9000"
    conn, err := net.Dial("tcp", seradd)
    if err != nil {
        fmt.Println(err)
    }

    res , err := conn.Write([]byte(" wanna join your server"))
    if err != nil {
        fmt.Println(err)
    }
	
    fmt.Println(res)
    conn.Close()
}