package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

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

	LineData struct {
		Lines Lines
	}
	// Lines ：セリフ配列構造体
	Lines []struct {
		Line  string
		Scene string
		Kind  string
	}
)

func Init(slackConfigBuf []byte) error {
	var config Config
	convertErr := ReadJSONOnStruct(slackConfigBuf, &config)
	if convertErr != nil {
		fmt.Println(convertErr)
		return convertErr
	}

	lineBuf, lineReadErr := ioutil.ReadFile("./data/lines.json")
	if lineReadErr != nil {
		fmt.Println(convertErr)
		return lineReadErr
	}

	var lineData LineData
	lineConvertErr := ReadJSONOnStruct(lineBuf, &lineData)
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
			rand.Seed(time.Now().UnixNano())
			say := lineData.Lines[rand.Intn(len(lineData.Lines)-1)].Line
			say = strings.Replace(say, "○○", "デスピサロ先輩", -1)
			rtm.SendMessage(rtm.NewOutgoingMessage(say, ev.Channel))
		}
	}
	return nil
}

// ReadJsonOnStruct ：コンフィグファイルの内容を構造体に格納する
func ReadJSONOnStruct(buf []byte, target interface{}) error {
	err := json.Unmarshal(buf, &target)
	if err != nil {
		return err
	}
	return nil
}
