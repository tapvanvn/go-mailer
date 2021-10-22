package entity

type PubsubConfig struct {
	Provider       string `json:"Provider"`
	ConnnectString string `json:"ConnnectString"`
	Topic          string `json:"Topic"`
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
	ChannelCapability int           `json:"ChannelCapability"`
	NumTemplater      int           `json:"NumTemplater"`
}
