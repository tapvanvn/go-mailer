package entity

type HealthConfig struct {
	HTTP     bool `json:"HTTP"`
	HTTPPort int  `json:"HTTPPort"`
}
type PubsubConfig struct {
	Provider      string `json:"Provider"`
	ConnectString string `json:"ConnectString"`
	Topic         string `json:"Topic"`
	Response      bool   `json:"Response"`
	ResponseTopic string `json:"ResponseTopic"`
}

type SMTPConfig struct {
	Account  string `json:"Account"`
	Password string `json:"Password"`
	Server   string `json:"Server"`
	Port     string `json:"Port"`
}

type Config struct {
	Pubsub            *PubsubConfig `json:"Pubsub"`
	SMTP              *SMTPConfig   `json:"SMTP"`
	HealthStatus      *HealthConfig `json:"HealthStatus"`
	ChannelCapability int           `json:"ChannelCapability"`
	NumTemplater      int           `json:"NumTemplater"`
}
