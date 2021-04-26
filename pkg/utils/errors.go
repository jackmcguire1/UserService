package utils

import "fmt"

var (
	ErrNotFound   = fmt.Errorf("could not find item")
	AlreadyExists = fmt.Errorf("item already exists")
)
