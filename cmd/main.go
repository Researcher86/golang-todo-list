package main

import (
	"encoding/json"
	"io/ioutil"
	app "todo/internal"
)

func main() {
	file, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		panic(err)
	}

	config := app.Config{}
	if err := json.Unmarshal(file, &config); err != nil {
		panic(err)
	}

	server := app.NewServer(&config)
	err = server.Start()
	if err != nil {
		panic(err)
	}
}
