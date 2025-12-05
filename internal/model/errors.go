package model

import "errors"

var (
	ErrNoLinks     = errors.New("no links in request")
	ErrNoNums      = errors.New("no nums in request")
	ErrNumLessZero = errors.New("num should be greater than 0")
	ErrNumNotFound = errors.New("num is not found")
)
