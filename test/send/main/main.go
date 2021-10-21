package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tapvanvn/gomailer/system"
)

var configPath string = ""

func main() {
	if len(os.Args) != 2 {
		fmt.Println("test.sh <configPath>")
		os.Exit(1)
	}
	configPath = os.Args[1]
	_, err := system.Init(configPath, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = system.EmailServer.SendEmail(system.Config.SMTP.Account, "tapvanvn@yahoo.com", "test email", "test email")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Success")
}
