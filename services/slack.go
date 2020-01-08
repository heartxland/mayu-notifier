package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/nlopes/slack"
)

type (
	// Config ：コンフィグ構造体
	Config struct {
		SlackParams SlackParams
	}

	// SlackParams ：アクセストークン構造体
	SlackParams struct {
		UserAccessToken string
		BotAccessToken  string
	}
	// Lines ：セリフ配列構造体
	Lines struct {
		LineData LineData
	}

	// LineData ：セリフ詳細
	LineData struct {
		Line  string
		Scene string
		Kind  string
	}
)

func Init(slackConfigBuf []byte) error {
	var config Config
	convertErr := ReadJsonOnStruct(slackConfigBuf, &config)
	if convertErr != nil {
		fmt.Println(convertErr)
		return convertErr
	}

	lineBuf, lineReadErr := ioutil.ReadFile("./data/lines.json")
	if lineReadErr != nil {
		fmt.Println(convertErr)
		return lineReadErr
	}

	var lines Lines
	lineConvertErr := ReadJsonOnStruct(lineBuf, &lines)
	if lineConvertErr != nil {
		fmt.Println(lineConvertErr)
		return lineConvertErr
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
	return nil
}

// ReadJsonOnStruct ：コンフィグファイルの内容を構造体に格納する
func ReadJsonOnStruct(buf []byte, target interface{}) error {
	err := json.Unmarshal(buf, &target)
	if err != nil {
		return err
	}
	return nil
}
