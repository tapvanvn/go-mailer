package entity

type SendRequest struct {
	RequestID    string            `json:"RequestID"`
	EmailAddress string            `json:"EmailAddress"`
	Params       map[string]string `json:"Params"`
}
