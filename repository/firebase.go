package repository

type FireBase interface {
	InsertEntry(message string) (*int, error)
}

type fireBase struct {
}

func NewFireBase() FireBase {
	return &fireBase{}
}

func (f *fireBase) InsertEntry(message string) (*int, error) {
	//Should insert generated insult into firebase
	//Should allow error to bubble up upon failure
	return nil, nil
}
