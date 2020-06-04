package model

type InsultStat struct {
	Verbs      []VerbCount      `json:"verbs"`
	Adjectives []AdjectiveCount `json:"adjectives"`
	Nouns      []NounCount      `json:"nouns"`
}

type VerbCount struct {
	Verb  string `json:"verb"`
	Count int    `json:"count"`
}

type AdjectiveCount struct {
	Adjective string `json:"adjective"`
	Count     int    `json:"count"`
}

type NounCount struct {
	Noun  string `json:"noun"`
	Count int    `json:"count"`
}
