package model

type Users struct {
	To          string  `json:"to"`
	From        string  `json:"from"`
	PhoneNumber *string `json:"phoneNumber"`
}
