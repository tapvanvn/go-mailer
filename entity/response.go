package entity

type Response struct {
	RequestID    string `json:"RequestID"`
	EmailAddress string `json:"EmailAddress"`
	Success      bool   `json:"Success"`
}
