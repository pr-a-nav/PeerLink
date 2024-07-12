package main

import "fmt"

func meet(){
	fmt.Println("gonna start")
}

type Ipad struct{
    ipad []string
	port []int
}

func NewIpad(ip string, port int) Ipad {
    return Ipad{
        ipad:   []string{ip},  // Initialize IP slice with the provided value
        port: []int{port}, // Initialize Port slice with the provided value
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