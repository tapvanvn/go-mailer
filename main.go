package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/tapvanvn/gomailer/entity"
	"github.com/tapvanvn/gomailer/system"
)

var requestChannel chan *entity.SendRequest = nil

func onMessage(message string) {
	request := &entity.SendRequest{}
	err := json.Unmarshal([]byte(message), request)
	if err != nil {
		log.Println(err)
		return
	}
}

func process() {
	for {
		request := <-requestChannel
		_ = request
		err := system.EmailServer.SendEmail(system.Config.SMTP.Account, request.EmailAddress, "test email", "test email content")
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	configPath := ""
	if len(os.Args) == 2 {
		configPath = strings.TrimSpace(os.Args[1])
	}
	if len(configPath) == 0 {
		log.Fatal("configPath is empty")
	}

	chn, err := system.Init(configPath, onMessage)
	if err != nil {
		log.Fatal(err)
	}
	requestChannel = chn
	if err != nil {
		log.Fatal(err)
	}
	process()
}
