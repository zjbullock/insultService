package model

type UserInfo struct {
	Name       string `firestore:"username"`
	Experience int    `firestore:"exp"`
	Rank       string `firestore:"rank"`
	DocID      *string
}
