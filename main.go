package main

import (
	"flag"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}
	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

func main() {
	wslName := flag.String("run", "", "wsl2 name")
	flag.Parse()

	/*
		if *wslName != "" {
			if err := exec.Command("wsl", "-d", *wslName, "-u", "root", "/etc/init.wsl").Run(); err != nil {
				panic(err)
			}
		}
		time.Sleep(5 * time.Second)
	*/

	// 获取WSL2 IP地址
	result, err := exec.Command("wsl", "ifconfig", "eth0").Output()
	if err != nil {
		panic(err)
	}

	var valid = regexp.MustCompile("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}")
	wslIP := valid.FindAllString(string(result), -1)[0]

	res, err := exec.Command("wsl", "netstat", "-lntp").Output()
	if err != nil {
		panic(err)
	}
	var reg = regexp.MustCompile(":(\\d{1,5})")
	// 去重
	ports := removeDuplicates(reg.FindAllString(string(res), -1))

	// netsh interface portproxy add v4tov4 listenport=80 connectaddress=172.17.83.208 connectport=80 listenaddress=* protocol=tcp
	for _, port := range ports {
		port = strings.TrimPrefix(port, ":")
		if err := exec.Command("netsh", "interface", "portproxy", "add", "v4tov4", fmt.Sprintf("listenport=%s", port), fmt.Sprintf("connectaddress=%s", wslIP), fmt.Sprintf("connectport=%s", port), "listenaddress=*", "protocol=tcp").Run(); err != nil {
			panic(err)
		}
	}
}
