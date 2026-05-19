package morphology

import (
	"iuno-api/models"
	"strings"
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
		case 3: 
			return generateThirdDeclension(lemma)
		case 4:
			return generateFourthDeclension(lemma)
		case 5: 
			return generateFifthDeclension(lemma)
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

func generateThirdDeclension(lemma models.Lemma) []models.Form {

	stem := removeEnding(*lemma.Genitive, "is")
	var endings map[string]map[string]string

	if (*lemma.Gender == "neuter") {
		endings = map[string]map[string]string{
			"singular": {
				"genitive":   "is",
				"dative":     "ī",
				"accusative": "us",
				"ablative":   "e",
				"vocative":   "us",
			},
			"plural": {
				"nominative": "a",
				"genitive":   "um",
				"dative":     "ibus",
				"accusative": "a",
				"ablative":   "ibus",
				"vocative":   "a",
			},
		}
	}else {
		endings = map[string]map[string]string{
			"singular": {
				"genitive":   "is",
				"dative":     "ī",
				"accusative": "em",
				"ablative":   "e",
				// "vocative":   "e",
			},
			"plural": {
				"nominative": "ēs",
				"genitive":   "um",
				"dative":     "ibus",
				"accusative": "ēs",
				"ablative":   "ibus",
				// "vocative":   "ī",
			},
		}
	}

	return buildNounForms(lemma, stem, endings)
}

func generateFourthDeclension(lemma models.Lemma) []models.Form {

	stem := removeEnding(*lemma.Genitive, "ūs")
	var endings map[string]map[string]string

	if (*lemma.Gender == "neuter") {
		endings = map[string]map[string]string{
			"singular": {
				"genitive":   "ūs",
				"dative":     "ū",
				"accusative": "ū",
				"ablative":   "ū",
				// "vocative":   "us",
			},
			"plural": {
				"nominative": "ua",
				"genitive":   "uum",
				"dative":     "ibus",
				"accusative": "ua",
				"ablative":   "ibus",
				// "vocative":   "a",
			},
		}
	}else {
		endings = map[string]map[string]string{
			"singular": {
				"genitive":   "ūs",
				"dative":     "uī",
				"accusative": "um",
				"ablative":   "ū",
				// "vocative":   "us",
			},
			"plural": {
				"nominative": "ūs",
				"genitive":   "uum",
				"dative":     "ibus",
				"accusative": "ūs",
				"ablative":   "ibus",
				// "vocative":   "a",
			},
		}
	}

	return buildNounForms(lemma, stem, endings)
}

func generateFifthDeclension(lemma models.Lemma) []models.Form {

	var stem string

	var endings map[string]map[string]string

	if strings.HasSuffix(*lemma.Genitive, "ēī") {
		stem = removeEnding(*lemma.Genitive, "ēī")
		endings = map[string]map[string]string{
		"singular": {
			"genitive":   "ēī",
			"dative":     "ēī",
			"accusative": "em",
			"ablative":   "ē",
		},
		"plural": {
			"nominative": "ēs",
			"genitive":   "ērum",
			"dative":     "ēbus",
			"accusative": "ēs",
			"ablative":   "ēbus",
		},
	}
	} else {
		stem = removeEnding(*lemma.Genitive, "eī")
		endings = map[string]map[string]string{
		"singular": {
			"genitive":   "eī",
			"dative":     "eī",
			"accusative": "em",
			"ablative":   "ē",
		},
		"plural": {
			"nominative": "ēs",
			"genitive":   "ērum",
			"dative":     "ēbus",
			"accusative": "ēs",
			"ablative":   "ēbus",
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

			// 3rd declension
			if *lemma.Declension == 3 {
				// neuter sing. nom == acc == voc
				if *lemma.Gender == "neuter" {
					if number == "singular" {
						if c == "accusative" {
							form = lemma.Lemma
						} else if c == "vocative" {
							form = lemma.Lemma
						}
					}
				}
				// sing. voc. == sing. nom.
				if number == "singular" && c == "vocative" {
					form = lemma.Lemma
				}
				// I-Stem
				if number == "plural" && c == "genitive" {
					if strings.HasSuffix(lemma.Lemma, "is") || 
						(*lemma.Gender == "neuter" && 
						(strings.HasSuffix(lemma.Lemma, "e") ||
						strings.HasSuffix(lemma.Lemma, "al") ||
						strings.HasSuffix(lemma.Lemma, "ar"))){
							form = stem + "ium"
					}
				}

				// I-Stem Neuter Plural
				if *lemma.Gender == "neuter" &&
					number == "plural" &&
					(c == "nominative" ||
					c == "accusative" ||
					c == "vocative") {

					if strings.HasSuffix(lemma.Lemma, "e") ||
						strings.HasSuffix(lemma.Lemma, "al") ||
						strings.HasSuffix(lemma.Lemma, "ar") {

						form = stem + "ia"
					}
				}
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