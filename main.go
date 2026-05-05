package main

import "log"

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	client := NewClient(config.Nick, config.User, config.Server, config.Port)
	if err = client.connect(); err != nil {
		log.Fatal(err)
	}
	defer client.conn.Close()
	client.register()
	client.run()
}
