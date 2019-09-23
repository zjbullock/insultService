package model

type Insult struct {
	Message *string `json:"message"`
	Id      *string `json:"fireStoreId"`
}

type InsultContent struct {
	Verb      string `json:"verb"`
	Adjective string `json:"adjective"`
	Noun      string `json:"noun"`
}
