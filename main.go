package main

import (
	"fmt"
	//"os/exec"
)

var (
	c   chan string //lock of getHost()
	str string      //result of getHost()
)

func main() {
	str := make([]string, 10)
	fmt.Println("process begins")
	getDomainList(str)
	fmt.Println(str)
}
