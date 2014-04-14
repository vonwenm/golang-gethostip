package main

import (
	//"fmt"
	"os/exec"
	"strings"
)

func getHost(domain string) {
	cmd := exec.Command("gethostip", "-d", domain)
	x, err := cmd.Output()
	if err != nil {
		lock <- false
		return
	}

	//
	res := string(x)
	domains[domain] = strings.Trim(res, "\r\n")
	lock <- true
	return
}
