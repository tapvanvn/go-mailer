package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tapvanvn/go-mailer/email"
	"github.com/tapvanvn/go-mailer/entity"
	"github.com/tapvanvn/go-mailer/system"
)

var configPath string = ""

func main() {

	rootPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(workDir)
	rootPath = workDir

	if len(os.Args) != 2 {
		fmt.Println("test.sh <configPath>")
		os.Exit(1)
	}
	configPath = os.Args[1]
	_, err = system.Init(rootPath+"/../..", configPath, nil)
	if err != nil {
		log.Fatal(err)
	}
	request := &entity.SendRequest{
		EmailAddress: "tapvanvn@yahoo.com",
		Params: map[string]map[string]string{
			"user": {
				"user_name": "Nguyễn Duy",
			},
		},
		Title:    "This is a test email",
		Template: "invite.html",
	}
	email.Send(request)

}
