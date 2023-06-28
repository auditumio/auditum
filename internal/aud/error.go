package aud

import "errors"

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrRecordNotFound  = errors.New("record not found")

	ErrDisabled = errors.New("disabled")
)
