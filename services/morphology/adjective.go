package morphology

import "iuno-api/models"

//
// ADJECTIVE MORPHOLOGY ENGINE
//
// This file generates:
//
// - masculine forms
// - feminine forms
// - neuter forms
//
// across:
//
// - singular
// - plural
// - all grammatical cases
//
// Supported:
//
// - 1st/2nd declension adjectives
//
// Example:
//
// bonus, bona, bonum
//

func GenerateAdjective(
	word models.Word,
) []models.Form {

	switch word.Declension {

	case 1:
		return generateFirstSecondDeclensionAdjective(word)
	}

	return []models.Form{}
}

// =====================================================
// 1ST / 2ND DECLENSION ADJECTIVES
// =====================================================
//
// bonus, bona, bonum
//
// stem = bon-
//

func generateFirstSecondDeclensionAdjective(
	word models.Word,
) []models.Form {

	stem := word.Stem

	if stem == "" {
		stem = removeEnding(word.Lemma, "us")
	}

	var forms []models.Form

	forms = append(
		forms,
		buildMasculineAdjectiveForms(word, stem)...,
	)

	forms = append(
		forms,
		buildFeminineAdjectiveForms(word, stem)...,
	)

	forms = append(
		forms,
		buildNeuterAdjectiveForms(word, stem)...,
	)

	return forms
}

// =====================================================
// MASCULINE
// =====================================================

func buildMasculineAdjectiveForms(
	word models.Word,
	stem string,
) []models.Form {

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

	return buildAdjectiveForms(
		stem,
		"masculine",
		endings,
	)
}

// =====================================================
// FEMININE
// =====================================================

func buildFeminineAdjectiveForms(
	word models.Word,
	stem string,
) []models.Form {

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

	return buildAdjectiveForms(
		stem,
		"feminine",
		endings,
	)
}

// =====================================================
// NEUTER
// =====================================================

func buildNeuterAdjectiveForms(
	word models.Word,
	stem string,
) []models.Form {

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

	return buildAdjectiveForms(
		stem,
		"neuter",
		endings,
	)
}

// =====================================================
// SHARED ADJECTIVE BUILDER
// =====================================================

func buildAdjectiveForms(
	stem string,
	gender string,
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

			forms = append(forms, models.Form{
				Form: stem + endings[number][grammaticalCase],

				Part: "adjective",

				Case: grammaticalCase,

				Number: number,

				Gender: gender,
			})
		}
	}

	return forms
}