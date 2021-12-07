package apihandler

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func isPlusSize(value string) bool {
	for _, size := range []string{"G1", "G2", "G3", "G4", "G5", "50", "52", "54"} {
		if strings.Contains(value, size) {
			return true
		}
	}
	return false
}

func mn(value string) string {
	var c conv
	t := transform.Chain(norm.NFD, runes.Remove(c), norm.NFC)
	result, _, _ := transform.String(t, value)

	return strings.ToLower(result)
}

//Title ...
func Title(value string) string {
	return strings.Title(strings.ToLower(value))
}

type conv struct{}

func (c conv) Contains(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

type class struct {
	Name         string
	Departamento string
	Categoria    string
}
