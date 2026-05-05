package main

import "log"

func main() {
	client := NewClient("lpall", "Liam Pallett", "irc.libera.chat", 6697)
	if err := client.connect(); err != nil {
		log.Fatal(err)
	}
	defer client.conn.Close()
	client.register()
	client.run()
}
