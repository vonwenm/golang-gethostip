package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func getDomainList(domains []string) error {
	bytes, err := ioutil.ReadFile(src_file_name)
	if nil != err {
		fmt.Println(err)
		return err
	}

	domains = strings.Split(string(bytes), "\n")
	return nil
}
