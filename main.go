package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func run(command string) string {
	args := strings.Split(command, " ")
	cmd := exec.Command(args[0], args[1:]...)

	var out bytes.Buffer
	cmd.Stdout = &out

	cmd.Run()

	return out.String()
}

func main() {
	fmt.Printf("hi")
}
