package system

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/tapvanvn/go-mailer/entity"
	"github.com/tapvanvn/gopubsubengine"
	"github.com/tapvanvn/gopubsubengine/awssqs"
	"github.com/tapvanvn/gopubsubengine/gpubsub"
	"github.com/tapvanvn/gopubsubengine/wspubsub"
	"github.com/tapvanvn/gotemplater"
	"github.com/tapvanvn/goutil"
)

var subscriber gopubsubengine.Subscriber = nil
var Publisher gopubsubengine.Publisher = nil

var HtmlRuntime = gotemplater.CreateHTMLRuntime()

var EmailServer *goutil.SmtpServer = nil
var Config *entity.Config = nil

var ErrBadConnectionString = errors.New("Bad Connectionstring")
var ErrEstablishConnection = errors.New("Cannot Establish connection")

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

		fmt.Printf("pubsub connect string:%s\n", Config.Pubsub.ConnectionString)

		topic := Config.Pubsub.Topic
		resTopic := Config.Pubsub.ResponseTopic

		if Config.Pubsub.Provider == "wspubsub" {
			newHub, err := wspubsub.NewWSPubSubHub(Config.Pubsub.ConnectionString)
			if err != nil {
				return nil, err
			}
			hub = newHub
		} else if Config.Pubsub.Provider == "gpubsub" {
			newHub, err := gpubsub.NewGPubSubHub(Config.Pubsub.ConnectionString)
			if err != nil {
				return nil, err
			}
			hub = newHub
		} else if Config.Pubsub.Provider == "awssqs" {
			parts := strings.Split(Config.Pubsub.ConnectionString, ":")
			if len(parts) != 3 {

				return nil, ErrBadConnectionString
			}
			sessionConfig := &aws.Config{
				Region:      aws.String(parts[0]),
				Credentials: credentials.NewStaticCredentials(parts[1], parts[2], ""),
			}

			sess, err := session.NewSession(sessionConfig)
			if err != nil {
				return nil, ErrEstablishConnection
			}
			newHub, err := awssqs.NewAWSSQSHubFromSession(sess)
			if err != nil {
				return nil, err
			}
			topicParts := strings.Split(Config.Pubsub.Topic, "@")
			if len(topicParts) != 2 {

				return nil, ErrBadConnectionString
			}
			topic = topicParts[0]
			newHub.SetTopicQueueURL(topicParts[0], topicParts[1])
			if len(Config.Pubsub.ResponseTopic) > 0 {
				topicParts = strings.Split(Config.Pubsub.ResponseTopic, "@")
				if len(topicParts) != 2 {

					return nil, ErrBadConnectionString
				}
				newHub.SetTopicQueueURL(topicParts[0], topicParts[1])
				resTopic = topicParts[0]
			}
			hub = newHub
		}
		if hub == nil {
			return nil, errors.New("Config error")
		}
		sub, err := hub.SubscribeOn(topic)
		if err != nil {
			return nil, err
		}
		subscriber = sub
		subscriber.SetProcessor(processMessage)
		if Config.Pubsub.Response {
			if len(resTopic) == 0 {
				return nil, errors.New("ReponseTopic cannot be empty")
			}
			pub, err := hub.PublishOn(resTopic)
			if err != nil {
				return nil, err
			}
			Publisher = pub
		}
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
