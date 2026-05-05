package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Message struct {
	prefix     string
	command    string
	parameters []string
}

func parse(line string) Message {
	var prefix string
	var command string
	var parameters []string

	if line[0] == ':' {
		prefix = line[1:strings.Index(line, " ")]
		line = strings.SplitN(line, " ", 2)[1]
	}

	command = line[0:strings.Index(line, " ")]
	line = strings.SplitN(line, " ", 2)[1]

	for line != "" {
		if line[0] == ':' {
			parameters = append(parameters, line[1:])
			break
		}
		index := strings.Index(line, " ")
		if index == -1 {
			parameters = append(parameters, line)
			break
		}
		parameters = append(parameters, line[0:index])
		line = strings.SplitN(line, " ", 2)[1]
	}

	return Message{prefix, command, parameters}
}

func (msg Message) String() string {
	var builder strings.Builder

	if msg.prefix != "" {
		builder.WriteByte(':')
		builder.WriteString(msg.prefix)
		builder.WriteByte(' ')
	}

	builder.WriteString(msg.command)

	if len(msg.parameters) > 0 {
		for _, element := range msg.parameters[:len(msg.parameters)-1] {
			builder.WriteByte(' ')
			builder.WriteString(element)
		}

		builder.WriteString(" :")
		builder.WriteString(msg.parameters[len(msg.parameters)-1])
	}

	builder.WriteString("\r\n")

	return builder.String()
}

func main() {
	conn, err := net.Dial("tcp", "irc.libera.chat:6667")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	nick := Message{"", "NICK", []string{"lpall"}}
	fmt.Fprintf(conn, "%s", nick)

	user := Message{"", "USER", []string{"lpall", "0", "*", "Liam Pallett"}}
	fmt.Fprintf(conn, "%s", user)

	buffServer := make(chan string)
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			msg := scanner.Text()
			buffServer <- msg
		}
		close(buffServer)
	}()

	buffClient := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := scanner.Text()
			buffClient <- msg
		}
	}()

	for {
		select {
		case line, ok := <-buffServer:
			if !ok {
				return
			}
			msg := parse(line)
			fmt.Println(line)
			switch msg.command {
			case "PING":
				pong := Message{"", "PONG", msg.parameters}
				fmt.Fprintf(conn, "%s", pong)
			}
		case line := <-buffClient:
			fmt.Fprintf(conn, "%s\r\n", line)
		}
	}
}
