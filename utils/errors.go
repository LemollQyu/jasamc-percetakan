package utils

import "errors"

var (
	ErrNameExists = errors.New("name category sudah digunakan")
	ErrSlugExists = errors.New("slug category sudah digunakan")
)
