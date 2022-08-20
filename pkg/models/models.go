package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created string
	Expires time.Time
}
