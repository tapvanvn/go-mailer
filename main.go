package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/tapvanvn/go-mailer/email"
	"github.com/tapvanvn/go-mailer/entity"
	"github.com/tapvanvn/go-mailer/system"
)

var requestChannel chan *entity.SendRequest = nil

func onMessage(message string) {
	//fmt.Println(message)
	request := &entity.SendRequest{}
	err := json.Unmarshal([]byte(message), request)
	if err != nil {
		log.Println(err)
		return
	}
	requestChannel <- request
}

func send(request *entity.SendRequest) {
	err := email.Send(request)
	if err != nil {
		log.Println(err)
		if system.Publisher != nil {
			res := &entity.Response{
				RequestID:    request.RequestID,
				EmailAddress: request.EmailAddress,
				Success:      false,
			}
			go system.Publisher.Publish(res)
		}
		return
	}
	if system.Publisher != nil {
		res := &entity.Response{
			RequestID:    request.RequestID,
			EmailAddress: request.EmailAddress,
			Success:      true,
		}
		go system.Publisher.Publish(res)
	}
}

func process() {
	for {
		request := <-requestChannel
		go send(request)
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

	configPath := "config/config.jsonc"
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
	if system.Config.HealthStatus != nil && system.Config.HealthStatus.HTTP {
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("i am ok"))
		})
		port := 80
		if system.Config.HealthStatus.HTTPPort > 0 {
			port = system.Config.HealthStatus.HTTPPort
		}
		go process()
		fmt.Printf("healthz on %d\n", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	} else {
		process()
	}
}
