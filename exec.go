package main

import (
	"os/exec"
	"strings"
)

func getHost(domain string) {
	cmd := exec.Command("gethostip", "-d", domain)
	x, err := cmd.Output()
	if err == nil {
		lock <- false
	}
	res := string(x)
	domains[domain] = strings.Trim(res, "\r\n")
	lock <- true
	return
}
