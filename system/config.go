package system

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/tapvanvn/gomailer/entity"
	"github.com/tapvanvn/gopubsubengine"
	"github.com/tapvanvn/gopubsubengine/gpubsub"
	"github.com/tapvanvn/gopubsubengine/wspubsub"
	"github.com/tapvanvn/gotemplater"
	"github.com/tapvanvn/goutil"
)

var subscriber gopubsubengine.Subscriber = nil
var HtmlRuntime = gotemplater.CreateHTMLRuntime()

var EmailServer *goutil.SmtpServer = nil
var Config *entity.Config = nil

func Init(rootPath string, configPath string, processMessage func(string)) (chan *entity.SendRequest, error) {
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}

	configData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	configData = goutil.TripJSONComment(configData)
	Config = &entity.Config{}
	err = json.Unmarshal(configData, Config)
	if err != nil {
		log.Fatal(err)
	}
	if Config.Pubsub != nil {
		var hub gopubsubengine.Hub = nil
		if Config.Pubsub.Provider == "wspubsub" {
			newHub, err := wspubsub.NewWSPubSubHub(Config.Pubsub.ConnnectString)
			if err != nil {
				return nil, err
			}
			hub = newHub
		} else if Config.Pubsub.Provider == "gpubsub" {
			newHub, err := gpubsub.NewGPubSubHub(Config.Pubsub.ConnnectString)
			if err != nil {
				return nil, err
			}
			hub = newHub
		}
		if hub == nil {
			return nil, errors.New("Config error")
		}
		sub, err := hub.SubscribeOn(Config.Pubsub.Topic)
		if err != nil {
			return nil, err
		}
		subscriber = sub
		subscriber.SetProcessor(processMessage)
	}

	EmailServer = goutil.NewSmtpServer(Config.SMTP.Server, Config.SMTP.Port, Config.SMTP.Account, Config.SMTP.Password)

	if EmailServer == nil {

		return nil, errors.New("Create Email Server fail")
	}

	err = gotemplater.InitTemplater(Config.NumTemplater)
	if err != nil {

		return nil, err
	}
	err = gotemplater.Templater.AddNamespace("page", rootPath+"/template")
	if err != nil {
		log.Fatal(err)
	}
	return make(chan *entity.SendRequest, Config.ChannelCapability), nil
}
