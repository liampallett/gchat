package main

import (
	"errors"
	"fmt"
	"strings"
)

func (client *Client) cmdQuit(args string) (Message, error) {
	if args != "" {
		return Message{"", "QUIT", []string{args}}, nil
	}
	return Message{"", "QUIT", []string{}}, nil
}

func (client *Client) cmdNick(args string) (Message, error) {
	if args != "" {
		return Message{"", "NICK", []string{args}}, nil
	}
	return Message{}, errors.New("need new nick")
}

func (client *Client) cmdJoin(args string) (Message, error) {
	if args != "" {
		client.currentChannel = args
		return Message{"", "JOIN", []string{client.currentChannel}}, nil
	}
	return Message{}, errors.New("specify channel to join")
}

func (client *Client) cmdMsg(args string) (Message, error) {
	parts := strings.SplitN(args, " ", 2)
	if len(parts) > 1 {
		return Message{"", "PRIVMSG", []string{parts[0], parts[1]}}, nil
	}
	return Message{}, errors.New("specify nick and message")
}

func (client *Client) cmdPart(args string) (Message, error) {
	channel := client.currentChannel
	partMsg := ""

	if args != "" {
		if args[0] == '#' {
			arg := strings.SplitN(args, " ", 2)
			channel = arg[0]
			if len(arg) > 1 {
				partMsg = arg[1]
			}
		} else {
			partMsg = args
		}
	}

	if channel == "" {
		return Message{}, errors.New("not in a channel")
	}
	if channel == client.currentChannel {
		client.currentChannel = ""
	}
	return Message{"", "PART", []string{channel, partMsg}}, nil
}

func (client *Client) cmdMe(args string) (Message, error) {
	if args != "" {
		if client.currentChannel == "" {
			return Message{}, errors.New("you are not in a channel")
		}
		return Message{"", "PRIVMSG", []string{client.currentChannel, "\x01ACTION " + args + "\x01"}}, nil
	}
	return Message{}, errors.New("specify action")
}

func (client *Client) cmdIgnore(args string) (Message, error) {
	if args != "" {
		client.ignored[args] = true
		return Message{}, nil
	}
	return Message{}, errors.New("specify nick to add to ignore list")
}

func (client *Client) cmdUnignore(args string) (Message, error) {
	if args != "" {
		delete(client.ignored, args)
		return Message{}, nil
	}
	return Message{}, errors.New("specify nick to remove from ignore list")
}

func (client *Client) cmdIgnores(args string) (Message, error) {
	if len(client.ignored) > 0 {
		for nick := range client.ignored {
			fmt.Printf("%s\n", nick)
		}
		return Message{}, nil
	} else {
		return Message{}, errors.New("no ignored users")
	}
}
