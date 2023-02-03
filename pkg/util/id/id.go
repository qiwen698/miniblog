package id

import (
	"strings"

	shortid "github.com/jasonsoft/go-short-id"
)

func GenShortID() string {
	opt := shortid.Options{
		Number:        4,
		StartWithYear: true,
		EndWithHost:   false,
	}
	return strings.ToLower(shortid.Generate(opt))
}
