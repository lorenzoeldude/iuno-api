package morphology

import "iuno-api/models"


func GenerateAdjective(lemma models.Lemma,) []models.Form {

	if lemma.Declension == nil {
		return []models.Form{}
	}

	switch *lemma.Declension {

	case 1:
		return generateFirstSecondDeclensionAdjective(lemma)
	}

	return []models.Form{}
}

// =====================================================
// 1ST / 2ND DECLENSION ADJECTIVES
// =====================================================

func generateFirstSecondDeclensionAdjective(
	lemma models.Lemma,
) []models.Form {

	stem := lemma.Lemma

	if stem == "" {
		stem = removeEnding(lemma.Lemma, "us")
	}

	var forms []models.Form

	forms = append(
		forms,
		buildMasculineAdjectiveForms(lemma, stem)...,
	)

	forms = append(
		forms,
		buildFeminineAdjectiveForms(lemma, stem)...,
	)

	forms = append(
		forms,
		buildNeuterAdjectiveForms(lemma, stem)...,
	)

	return forms
}

// =====================================================
// MASCULINE
// =====================================================

func buildMasculineAdjectiveForms(
	lemma models.Lemma,
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
	lemma models.Lemma,
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
	lemma models.Lemma,
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