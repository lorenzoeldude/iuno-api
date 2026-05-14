package morphology

import "iuno-api/models"

//
// NOUN MORPHOLOGY ENGINE (FULL SYSTEM)
// Supports:
// 1st, 2nd, 3rd, 4th, 5th declensions
//

func GenerateNoun(word models.Word) []models.Form {

	switch word.Declension {

	case 1:
		return generateFirstDeclension(word)

	case 2:
		if word.Gender == "neuter" {
			return generateSecondDeclensionNeuter(word)
		}
		return generateSecondDeclensionMasculine(word)

	case 3:
		return generateThirdDeclension(word)

	case 4:
		return generateFourthDeclension(word)

	case 5:
		return generateFifthDeclension(word)
	}

	return []models.Form{}
}

//
// =====================================================
// 1ST DECLENSION
// =====================================================
//

func generateFirstDeclension(word models.Word) []models.Form {

	stem := word.Stem
	if stem == "" {
		stem = removeEnding(word.Lemma, "a")
	}

	endings := map[string]map[string]string{
		"singular": {
			"nominative": "a",
			"genitive":   "ae",
			"dative":     "ae",
			"accusative": "am",
			"ablative":   "a",
			"vocative":   "a",
		},
		"plural": {
			"nominative": "ae",
			"genitive":   "arum",
			"dative":     "is",
			"accusative": "as",
			"ablative":   "is",
			"vocative":   "ae",
		},
	}

	return buildNounForms(word, stem, endings)
}

//
// =====================================================
// 2ND DECLENSION
// =====================================================
//

func generateSecondDeclensionMasculine(word models.Word) []models.Form {

	stem := word.Stem
	if stem == "" {
		stem = removeEnding(word.Lemma, "us")
	}

	endings := map[string]map[string]string{
		"singular": {
			"nominative": "us",
			"genitive":   "i",
			"dative":     "o",
			"accusative": "um",
			"ablative":   "o",
			"vocative":   "e",
		},
		"plural": {
			"nominative": "i",
			"genitive":   "orum",
			"dative":     "is",
			"accusative": "os",
			"ablative":   "is",
			"vocative":   "i",
		},
	}

	return buildNounForms(word, stem, endings)
}

func generateSecondDeclensionNeuter(word models.Word) []models.Form {

	stem := word.Stem
	if stem == "" {
		stem = removeEnding(word.Lemma, "um")
	}

	endings := map[string]map[string]string{
		"singular": {
			"nominative": "um",
			"genitive":   "i",
			"dative":     "o",
			"accusative": "um",
			"ablative":   "o",
			"vocative":   "um",
		},
		"plural": {
			"nominative": "a",
			"genitive":   "orum",
			"dative":     "is",
			"accusative": "a",
			"ablative":   "is",
			"vocative":   "a",
		},
	}

	return buildNounForms(word, stem, endings)
}

//
// =====================================================
// 3RD DECLENSION (CORE FIX)
// =====================================================
//

func generateThirdDeclension(word models.Word) []models.Form {

	stem := word.Stem
	if stem == "" {
		stem = word.Lemma // fallback (you should improve later)
	}

	endings := map[string]map[string]string{
		"singular": {
			"nominative": "",
			"genitive":   "is",
			"dative":     "i",
			"accusative": "em",
			"ablative":   "e",
			"vocative":   "",
		},
		"plural": {
			"nominative": "es",
			"genitive":   "um",
			"dative":     "ibus",
			"accusative": "es",
			"ablative":   "ibus",
			"vocative":   "es",
		},
	}

	return buildNounForms(word, stem, endings)
}

//
// =====================================================
// 4TH DECLENSION
// =====================================================
//

func generateFourthDeclension(word models.Word) []models.Form {

	stem := word.Stem
	if stem == "" {
		stem = removeEnding(word.Lemma, "us")
	}

	endings := map[string]map[string]string{
		"singular": {
			"nominative": "us",
			"genitive":   "us",
			"dative":     "ui",
			"accusative": "um",
			"ablative":   "u",
			"vocative":   "us",
		},
		"plural": {
			"nominative": "us",
			"genitive":   "uum",
			"dative":     "ibus",
			"accusative": "us",
			"ablative":   "ibus",
			"vocative":   "us",
		},
	}

	return buildNounForms(word, stem, endings)
}

//
// =====================================================
// 5TH DECLENSION
// =====================================================
//

func generateFifthDeclension(word models.Word) []models.Form {

	stem := word.Stem
	if stem == "" {
		stem = removeEnding(word.Lemma, "es")
	}

	endings := map[string]map[string]string{
		"singular": {
			"nominative": "es",
			"genitive":   "ei",
			"dative":     "ei",
			"accusative": "em",
			"ablative":   "e",
			"vocative":   "es",
		},
		"plural": {
			"nominative": "es",
			"genitive":   "erum",
			"dative":     "ebus",
			"accusative": "es",
			"ablative":   "ebus",
			"vocative":   "es",
		},
	}

	return buildNounForms(word, stem, endings)
}

//
// =====================================================
// SHARED BUILDER
// =====================================================
//

func buildNounForms(
	word models.Word,
	stem string,
	endings map[string]map[string]string,
) []models.Form {

	var forms []models.Form

	numbers := []string{"singular", "plural"}
	cases := []string{
		"nominative",
		"genitive",
		"dative",
		"accusative",
		"ablative",
		"vocative",
	}

	for _, number := range numbers {
		for _, c := range cases {

			forms = append(forms, models.Form{
				Form: stem + endings[number][c],
				Part: "noun",
				Case: c,
				Number: number,
				Gender: word.Gender,
			})
		}
	}

	return forms
}

//
// =====================================================
// HELPERS
// =====================================================
//

func removeEnding(word string, ending string) string {
	if len(word) < len(ending) {
		return word
	}
	return word[:len(word)-len(ending)]
}