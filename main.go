package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tapvanvn/gomailer/email"
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
		go email.Send(request)
	}
}

func main() {
	rootPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	} else {
		rootPath = workDir
	}

	configPath := ""
	if len(os.Args) == 2 {
		configPath = strings.TrimSpace(os.Args[1])
	}
	if len(configPath) == 0 {
		log.Fatal("configPath is empty")
	}

	chn, err := system.Init(rootPath, configPath, onMessage)
	if err != nil {
		log.Fatal(err)
	}
	requestChannel = chn
	if err != nil {
		log.Fatal(err)
	}
	process()
}
