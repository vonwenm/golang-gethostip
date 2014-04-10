package main

import (
	"os/exec"
)

func getHost(domain string) {
	cmd := exec.Command("gethostip", "-d", domain)
	x, _ := cmd.Output()
	res := string(x)
	c <- res
}
