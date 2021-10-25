package email

import (
	"fmt"

	"github.com/tapvanvn/go-mailer/entity"
	"github.com/tapvanvn/go-mailer/system"
	"github.com/tapvanvn/gosmartstring"
	"github.com/tapvanvn/gotemplater"
)

func Send(request *entity.SendRequest) error {

	renderContext := gosmartstring.CreateContext(system.HtmlRuntime)

	for objName, define := range request.Params {

		renderContext.RegisterInterface(objName, define)
	}

	content, err := gotemplater.Templater.Render(fmt.Sprintf("page:%s", request.Template), renderContext)
	if err != nil {

		return err
	}
	err = system.EmailServer.SendEmail(system.Config.SMTP.Account, request.EmailAddress, request.Title, content)
	if err != nil {

		return err
	}
	return nil
}
