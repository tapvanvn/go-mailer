package email

import (
	"fmt"
	"log"

	"github.com/tapvanvn/gomailer/entity"
	"github.com/tapvanvn/gomailer/system"
	"github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotemplater"
)

func Send(request *entity.SendRequest) {

	renderContext := gosmartstring.CreateContext(system.HtmlRuntime)

	for objName, define := range request.Params {

		renderContext.RegisterInterface(objName, define)
	}

	content, err := gotemplater.Templater.Render(fmt.Sprintf("page:%s", request.Template), renderContext)
	if err != nil {

		log.Println(err)
		return
	}
	err = system.EmailServer.SendEmail(system.Config.SMTP.Account, request.EmailAddress, request.Title, content)
	if err != nil {
		log.Println(err)
		return
	}
}
