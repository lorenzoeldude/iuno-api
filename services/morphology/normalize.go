package morphology

import "strings"

func NormalizeLatin(s string) string {
	replacer := strings.NewReplacer(
		"ā", "a",
		"ē", "e",
		"ī", "i",
		"ō", "o",
		"ū", "u",
		"ȳ", "y",

		"Ā", "a",
		"Ē", "e",
		"Ī", "i",
		"Ō", "o",
		"Ū", "u",
		"Ȳ", "y",
	)

	return strings.ToLower(replacer.Replace(s))
}