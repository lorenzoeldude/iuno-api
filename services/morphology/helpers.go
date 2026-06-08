package morphology

import (
	"iuno-api/models"
	"strings"
)

// HELPERS
func removeEnding(word string, ending string) string {
    if strings.HasSuffix(word, ending) {
        return strings.TrimSuffix(word, ending)
    }
    return word
}

func pppStem(supine string) string {
	return removeEnding(supine, "um") 
}

func removeVerbEnding(word string, ending string) string {
	if len(word) < len(ending) {
		return word
	}

	return word[:len(word)-len(ending)]
}

func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func buildGerundiveStem(lemma models.Lemma) string {

    infinitive := *lemma.Infinitive

    switch *lemma.Conjugation {

    case 1:
        return removeEnding(infinitive, "āre") + "and"

    case 2:
        return removeEnding(infinitive, "re") + "nd"

    case 3:
        return removeEnding(infinitive, "ere") + "end"

    case 31:
        return removeEnding(infinitive, "ere") + "iend"

    case 4:
        return removeEnding(infinitive, "re") + "end"
    }

    return ""
}

func verbForm(
	form string,
	person int,
	number string,
	tense string,
	mood string,
	voice string,
) models.Form {
	return models.Form{
		Form:         form,
		PartOfSpeech: "verb",

		Person: IntPtr(person),
		Number: number,

		Tense: StringPtr(tense),
		Mood:  StringPtr(mood),
		Voice: StringPtr(voice),
	}
}

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