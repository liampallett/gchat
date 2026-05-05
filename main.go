package main

import (
	"bufio"
	"crypto/tls"
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

type Client struct {
	nick           string
	user           string
	conn           net.Conn
	currentChannel string
}

func (client *Client) connect() {
	var err error
	client.conn, err = tls.Dial("tcp", "irc.libera.chat:6697", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (client *Client) register() {
	nick := Message{"", "NICK", []string{client.nick}}
	fmt.Fprintf(client.conn, "%s", nick)

	user := Message{"", "USER", []string{client.nick, "0", "*", client.user}}
	fmt.Fprintf(client.conn, "%s", user)
}

func (client *Client) run() {
	buffServer := make(chan string)
	go func() {
		scanner := bufio.NewScanner(client.conn)
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
				fmt.Fprintf(client.conn, "%s", pong)
			}
		case line := <-buffClient:
			if line[0] != '/' {
				privmsg := Message{"", "PRIVMSG", []string{client.currentChannel, line}}
				fmt.Fprintf(client.conn, "%s", privmsg)
			} else {
				rawCommand := line[1:strings.Index(line, " ")]
				switch rawCommand {
				case "join":
					client.currentChannel = line[strings.Index(line, " ")+1:]
					join := Message{"", "JOIN", []string{client.currentChannel}}
					fmt.Fprintf(client.conn, "%s", join)
				}
			}
		}
	}
}

func main() {
	client := Client{nick: "lpall", user: "Liam Pallett"}
	client.connect()
	defer client.conn.Close()
	client.register()
	client.run()
}
