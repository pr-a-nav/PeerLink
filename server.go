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
func vishwas() {
	fmt.Println("hello pranav")
}

