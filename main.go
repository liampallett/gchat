package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "irc.libera.chat:6667")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	nick := "lpall"
	user := "Liam Pallett"

	fmt.Fprintf(conn, "NICK %s\r\n", nick)
	fmt.Fprintf(conn, "USER %s 0 * :%s\r\n", nick, user)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println(message)

		if strings.HasPrefix(message, "PING") {
			fmt.Fprintf(conn, "%s\r\n", strings.Replace(message, "PING", "PONG", 1))
		}
	}
}
