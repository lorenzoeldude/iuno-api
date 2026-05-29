package morphology

import (
	"iuno-api/models"
)

func removeEnding(lemma string, ending string) string {
	if len(lemma) < len(ending) {
		return lemma
	}
	return lemma[:len(lemma)-len(ending)]
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