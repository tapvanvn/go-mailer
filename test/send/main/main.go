package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tapvanvn/go-mailer/system"
)

var configPath string = ""

func main() {
	if len(os.Args) != 2 {
		fmt.Println("test.sh <configPath>")
		os.Exit(1)
	}
	rootPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	} else {
		rootPath = workDir
	}

	configPath = os.Args[1]
	_, err = system.Init(rootPath+"/../..", configPath, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = system.EmailServer.SendEmail(system.Config.SMTP.Account, "tapvanvn@yahoo.com", "test email", "test email")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Success")
}
