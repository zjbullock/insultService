package model

type Titles struct {
	Titles []Title `firestore:"titles"`
}

type Title struct {
	Name       string `firestore:"name"`
	Experience int    `firestore:"exp"`
}
