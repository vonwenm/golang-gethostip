package main

import (
	"fmt"
	"io/ioutil"
	"strings"
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

func TestFile() {
	bytes, err := ioutil.ReadFile(src_file_name)
	if nil != err {
		fmt.Println(err)
		return err
	}

	temp := strings.Split(string(bytes), "\n")
	fmt.Println(temp)
}
