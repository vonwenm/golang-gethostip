package main

import (
	"fmt"
	"testing"
)

func TestChan(t *testing.T) {
	t.Log("testChan begins")
	var domains = []string{"gzidc.com", "gzidc.net", "cloverphp.com", "cloverphp.net"}

	c = make(chan string)
	for _, v := range domains {
		go getHost(v)
	}

	for i := 0; i < len(domains); i++ {
		str = <-c
		fmt.Print(str)
	}
	t.Log("testChan end")
}
