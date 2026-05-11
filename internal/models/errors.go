package models

import "errors"

var ErrDuplicateLink = errors.New("this link already exists in our database")

var ErrNotFound = errors.New("record not found")
