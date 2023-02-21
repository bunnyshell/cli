package common

type ItemOptions struct {
	Options

	ID string
}

func NewItemOptions(id string) *ItemOptions {
	return &ItemOptions{
		ID: id,
	}
}
