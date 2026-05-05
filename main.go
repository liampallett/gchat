package main

import "log"

func main() {
	client := Client{nick: "lpall", user: "Liam Pallett"}
	if err := client.connect(); err != nil {
		log.Fatal(err)
	}
	defer client.conn.Close()
	client.register()
	client.run()
}
