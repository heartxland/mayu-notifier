package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/nlopes/slack"
)

// Config ：コンフィグ構造体
type Config struct {
	SlackParams SlackParams
}

// SlackParams ：アクセストークン構造体
type SlackParams struct {
	UserAccessToken string
	BotAccessToken  string
}

// Lines ：セリフ配列構造体
type Lines struct {
	LineData LineData
}

// LineData ：セリフ詳細
type LineData struct {
	Line  string
	Scene string
	kind  string
}

func main() {
	buf, readErr := ioutil.ReadFile("./config/config.json")
	if readErr != nil {
		fmt.Println(readErr)
		return
	}

	var config Config
	convertErr := ReadConfigOnStruct(buf, &config)
	if convertErr != nil {
		fmt.Println(convertErr)
		return
	}

	api := slack.New(config.SlackParams.BotAccessToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			rtm.SendMessage(rtm.NewOutgoingMessage("まゆですよぉ", ev.Channel))
		}
	}
}

// ReadConfigOnStruct ：コンフィグファイルの内容を構造体に格納する
func ReadConfigOnStruct(buf []byte, config *Config) error {
	err := json.Unmarshal(buf, &config)
	if err != nil {
		return err
	}
	return nil
}
