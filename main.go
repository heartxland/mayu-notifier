package main

import (
	"fmt"
	"io/ioutil"

	"./services"
)

func main() {
	buf, readErr := ioutil.ReadFile("./config/config.json")
	if readErr != nil {
		fmt.Println(readErr)
		return
	}

	slackErr := services.Init(buf)
	if slackErr != nil {
		fmt.Println(slackErr)
		return
	}
}
