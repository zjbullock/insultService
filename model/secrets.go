package model

type Secrets struct {
	SID       string `json:"accountSid"`
	AuthToken string `json:"authToken"`
	Number    string `json:"number"`
}
