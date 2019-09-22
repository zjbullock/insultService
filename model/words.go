package model

type Words struct {
	Verb      []string `firestore:"verbs"`
	Adjective []string `firestore:"adjectives"`
	Noun      []string `firestore:"nouns"`
}
