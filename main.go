package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func am_i_the_admin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		fmt.Println("admin no")
		return false
	}
	fmt.Println("admin yes")
	return true
}

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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func start_the_wsl() {
	wslName := flag.String("run", "", "wsl2 name")
	flag.Parse()

	if *wslName != "" {
		if err := exec.Command("wsl", "-d", *wslName, "-u", "root", "/etc/init.wsl").Run(); err != nil {
			panic(err)
		}
	}
	time.Sleep(5 * time.Second)
}

func get_exists_port() []string {
	res, err := exec.Command("netsh", "interface", "portproxy", "show", "v4tov4").Output()
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%v", string(res))
	var reg = regexp.MustCompile("(\\d{1,5})\\s*\n")
	matchs := reg.FindAllStringSubmatch(string(res)+"\n", -1)
	var list = make([]string, 0)
	for _, l := range matchs {
		list = append(list, l[1])
	}
	ports := removeDuplicates(list)
	return ports
}

func make_the_port_forward() {
	// get WSL2 IP address
	result, err := exec.Command("wsl", "ifconfig", "eth0").Output()
	if err != nil {
		panic(err)
	}

	var valid = regexp.MustCompile("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}")
	wslIP := valid.FindAllString(string(result), -1)[0]
	fmt.Printf("\nwsl2 ip is: %v\n", string(wslIP))

	res, err := exec.Command("wsl", "netstat", "-lntp").Output()
	if err != nil {
		panic(err)
	}
	var reg = regexp.MustCompile(":(\\d{1,5})")
	ports := removeDuplicates(reg.FindAllString(string(res), -1))

	// netsh interface portproxy add v4tov4 listenport=80 connectaddress=172.17.83.208 connectport=80 listenaddress=* protocol=tcp
	exists_ports := get_exists_port()
	fmt.Printf("the exists_ports: %v\n", exists_ports)
	for _, port := range ports {
		port = strings.TrimPrefix(port, ":")
		if !contains(exists_ports, port) {
			fmt.Printf("proxyed port: %v\n", port)
			if err := exec.Command("netsh", "interface", "portproxy", "add", "v4tov4", fmt.Sprintf("listenport=%s", port), fmt.Sprintf("connectaddress=%s", wslIP), fmt.Sprintf("connectport=%s", port), "listenaddress=*", "protocol=tcp").Run(); err != nil {
				panic(err)
			}
		}
	}
}

func main() {
	if am_i_the_admin() == false {
		fmt.Printf("Run me as admin!")
		os.Exit(0)
	}

	for true {
		make_the_port_forward()
		time.Sleep(10 * time.Second)
	}
}
