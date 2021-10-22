package entity

type SendRequest struct {
	RequestID    string                       `json:"RequestID"`
	EmailAddress string                       `json:"EmailAddress"`
	Title        string                       `json:"Title"`
	Params       map[string]map[string]string `json:"Params"`
	Template     string                       `json:"Template"`
}
