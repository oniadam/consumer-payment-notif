package utils

import (
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func FormatRupiah(val string) string {
	p := message.NewPrinter(language.Indonesian)

	// parse float dulu
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return "0"
	}

	// buang desimal dengan convert ke int64
	n := int64(f)

	return p.Sprintf("%d", n)

}
