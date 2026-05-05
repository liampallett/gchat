package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	nick           string
	user           string
	server         string
	port           int
	conn           net.Conn
	currentChannel string
	handlers       map[string]func(Message)
}

func NewClient(nick, user, server string, port int) *Client {
	client := &Client{nick: nick, user: user, server: server, port: port}
	client.handlers = map[string]func(Message){
		"PING":    client.handlePing,
		"PRIVMSG": client.handlePrivmsg,
		"JOIN":    client.handleJoin,
		"PART":    client.handlePart,
		"QUIT":    client.handleQuit,
		"NICK":    client.handleNick,
	}
	return client
}

func (client *Client) connect() error {
	var err error
	client.conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", client.server, client.port), nil)
	return err
}

func (client *Client) register() {
	nick := Message{"", "NICK", []string{client.nick}}
	client.send(nick)

	user := Message{"", "USER", []string{client.nick, "0", "*", client.user}}
	client.send(user)
}

func (client *Client) send(msg Message) {
	fmt.Fprintf(client.conn, "%s", msg)
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
			if handler, ok := client.handlers[msg.command]; ok {
				handler(msg)
			} else {
				fmt.Println(line)
			}
		case line := <-buffClient:
			msg, err := client.parseInput(line)
			if err != nil {
				fmt.Println(err)
			} else {
				client.send(msg)
				if msg.command == "PRIVMSG" {
					fmt.Printf("<%s> %s\n", client.nick, msg.parameters[1])
				}
			}
		}
	}
}

func (client *Client) parseInput(line string) (Message, error) {
	if line[0] != '/' {
		if client.currentChannel == "" {
			return Message{}, errors.New("you are not in a channel")
		}
		return Message{"", "PRIVMSG", []string{client.currentChannel, line}}, nil
	}

	parts := strings.SplitN(line[1:], " ", 2)
	command := parts[0]
	args := ""
	if len(parts) > 1 {
		args = parts[1]
	}

	switch command {
	case "quit":
		return client.cmdQuit(args)
	case "nick":
		return client.cmdNick(args)
	case "join":
		return client.cmdJoin(args)
	case "msg":
		return client.cmdMsg(args)
	case "part":
		return client.cmdPart(args)
	default:
		return Message{}, errors.New("unrecognised command")
	}
}
