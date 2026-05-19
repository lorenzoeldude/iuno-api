package morphology

import (
	"iuno-api/models"
	// "strings"
)

//
// NOUN MORPHOLOGY ENGINE (FULL SYSTEM)
// Supports:
// 1st, 2nd, 3rd, 4th, 5th declensions
//

func GenerateNoun(lemma models.Lemma) []models.Form {

	switch *lemma.Declension {
		case 1:
			return generateFirstDeclension(lemma)

		case 2: 
			return generateSecondDeclension(lemma)
	}

	return []models.Form{}
}

//
// =====================================================
// 1ST DECLENSION
// =====================================================
//

func generateFirstDeclension(lemma models.Lemma) []models.Form {

	stem := removeEnding(*lemma.Genitive, "ae")

	endings := map[string]map[string]string{
		"singular": {
			"genitive":   "ae",
			"dative":     "ae",
			"accusative": "am",
			"ablative":   "ā",
			"vocative":   "a",
		},
		"plural": {
			"nominative": "ae",
			"genitive":   "ārum",
			"dative":     "īs",
			"accusative": "ās",
			"ablative":   "īs",
			"vocative":   "ae",
		},
	}

	return buildNounForms(lemma, stem, endings)
}

func generateSecondDeclension(lemma models.Lemma) []models.Form {

	stem := removeEnding(*lemma.Genitive, "ī")
	var endings map[string]map[string]string

	if (*lemma.Gender == "neuter") {
		endings = map[string]map[string]string{
		"singular": {
			"genitive":   "ī",
			"dative":     "ō",
			"accusative": "um",
			"ablative":   "ō",
			"vocative":   "um",
		},
		"plural": {
			"nominative": "a",
			"genitive":   "ōrum",
			"dative":     "īs",
			"accusative": "a",
			"ablative":   "īs",
			"vocative":   "a",
		},
	}
	} else {
		endings = map[string]map[string]string{
			"singular": {
				"genitive":   "ī",
				"dative":     "ō",
				"accusative": "um",
				"ablative":   "ō",
				"vocative":   "e",
			},
			"plural": {
				"nominative": "ī",
				"genitive":   "ōrum",
				"dative":     "īs",
				"accusative": "ōs",
				"ablative":   "īs",
				"vocative":   "ī",
			},
		}
	}

	return buildNounForms(lemma, stem, endings)
}

func buildNounForms(
	lemma models.Lemma,
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

			form := stem + endings[number][c]

			if number == "singular" && c == "nominative" {
				form = lemma.Lemma
			}

			forms = append(forms, models.Form{
				Form:   form,
				Part:   "noun",
				Case:   c,
				Number: number,
				Gender: *lemma.Gender,
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

func removeEnding(lemma string, ending string) string {
	if len(lemma) < len(ending) {
		return lemma
	}
	return lemma[:len(lemma)-len(ending)]
}