package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tapvanvn/go-mailer/system"
	"github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotemplater"
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

	htmlRuntime := gotemplater.CreateHTMLRuntime()

	renderContext := gosmartstring.CreateContext(htmlRuntime)

	gotemplater.Templater.Debug()

	_ = renderContext

	renderContext.RegisterInterface("user", map[string]string{
		"user_name": "Nguyễn Duy",
	})

	content, err := gotemplater.Templater.Render("page:invite.html", renderContext)

	if err != nil {

		log.Fatal(err)
	}
	fmt.Println(content)
	err = system.EmailServer.SendEmail(system.Config.SMTP.Account, "tapvanvn@yahoo.com", "test email", content)
	if err != nil {
		log.Fatal(err)
	}
}
