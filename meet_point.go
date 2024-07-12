package main

import "fmt"

func meet(){
	fmt.Println("gonna start")
}

type Ipad struct{
    ipad []string
	port []int
}

func (i Ipad )new (ipad string, port int ){
i.ipad = append(i.ipad,ipad )
i.port = append(i.port, port)

}