package morphology

import "iuno-api/models"

//
// NOUN MORPHOLOGY ENGINE
//
// This file generates:
//
// - singular forms
// - plural forms
// - all grammatical cases
//
// Supported:
//
// - 1st declension
// - 2nd declension masculine
// - 2nd declension neuter
//

func GenerateNoun(
	word models.Word,
) []models.Form {

	switch word.Declension {

	case 1:
		return generateFirstDeclension(word)

	case 2:

		// neuter nouns
		if word.Gender == "neuter" {
			return generateSecondDeclensionNeuter(word)
		}

		// masculine default
		return generateSecondDeclensionMasculine(word)
	}

	return []models.Form{}
}

// =====================================================
// 1ST DECLENSION
// =====================================================
//
// luna, lunae (f)
//
// stem = lun-
//

func generateFirstDeclension(
	word models.Word,
) []models.Form {

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

	return buildNounForms(
		word,
		stem,
		endings,
	)
}

// =====================================================
// 2ND DECLENSION MASCULINE
// =====================================================
//
// servus, servi
//
// stem = serv-
//

func generateSecondDeclensionMasculine(
	word models.Word,
) []models.Form {

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

	return buildNounForms(
		word,
		stem,
		endings,
	)
}

// =====================================================
// 2ND DECLENSION NEUTER
// =====================================================
//
// bellum, belli
//
// stem = bell-
//

func generateSecondDeclensionNeuter(
	word models.Word,
) []models.Form {

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

	return buildNounForms(
		word,
		stem,
		endings,
	)
}

// =====================================================
// SHARED NOUN BUILDER
// =====================================================

func buildNounForms(
	word models.Word,
	stem string,
	endings map[string]map[string]string,
) []models.Form {

	var forms []models.Form

	numbers := []string{
		"singular",
		"plural",
	}

	cases := []string{
		"nominative",
		"genitive",
		"dative",
		"accusative",
		"ablative",
		"vocative",
	}

	for _, number := range numbers {

		for _, grammaticalCase := range cases {

			form := models.Form{
				Form: stem + endings[number][grammaticalCase],

				Part: "noun",

				Case:   grammaticalCase,
				Number: number,

				Gender: word.Gender,
			}

			forms = append(forms, form)
		}
	}

	return forms
}

// =====================================================
// HELPERS
// =====================================================

func removeEnding(
	word string,
	ending string,
) string {

	if len(word) < len(ending) {
		return word
	}

	return word[:len(word)-len(ending)]
}