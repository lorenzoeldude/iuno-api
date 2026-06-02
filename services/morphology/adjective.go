package morphology

import (
	// "log"
	"iuno-api/models"
)

func GenerateAdjective(lemma models.Lemma) []models.Form {

	if lemma.Declension == nil {
		return []models.Form{}
	}

	switch *lemma.Declension {

	case 12:
		return generateFirstSecondDeclensionAdjective(lemma)

	case 31, 32, 33: 
		return generateThirdDeclensionAdjective(lemma)

	// case 32:
	// 	return generateThirdDeclensionTwoTerminationAdjective(lemma)

	// case 33:
	// 	return generateThirdDeclensionThreeTerminationAdjective(lemma)
	}

	return []models.Form{}
}

// =====================================================
// 1ST / 2ND DECLENSION ADJECTIVES
// =====================================================

func generateFirstSecondDeclensionAdjective(
	lemma models.Lemma,
) []models.Form {

	stem := removeEnding(*lemma.Genitive, "ī")

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
	
	forms = append(forms, buildComparativeForms(stem)...)

	forms = append(forms, buildSuperlativeForms(stem)...)

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

	return buildAdjectiveForms(
		stem,
		"neuter",
		endings,
	)
}

// =====================================================
// 3RD DECLENSION ADJECTIVES
// =====================================================

func generateThirdDeclensionAdjective(
	lemma models.Lemma,
) []models.Form {

	stem := removeEnding(*lemma.Genitive, "is")

	switch *lemma.Declension {
		case 31:
			return buildThirdDeclensionOneTerminationForms(lemma, stem)
		// case 32:
		// 	return buildThirdDeclensionTwoTerminationForms(lemma, stem)
		// case 33:
		// 	return buildThirdDeclensionThreeTerminationForms(lemma, stem)
	}

	return []models.Form{}
}

func buildThirdDeclensionOneTerminationForms(
	lemma models.Lemma,
	stem string,
) []models.Form {

	endings := map[string]map[string]string{

		"singular": {
			"genitive":   "is",
			"dative":     "ī",
			"accusative": "em",
			"ablative":   "ī",
		},

		"plural": {
			"nominative": "ēs",
			"genitive":   "ium",
			"dative":     "ibus",
			"accusative": "ēs",
			"ablative":   "ibus",
		},
	}

	var forms []models.Form

	forms = append(
		forms,
		buildThirdDeclensionForms(lemma, stem, "masculine", endings)...,
	)

	forms = append(
		forms,
		buildThirdDeclensionForms(lemma, stem, "feminine", endings)...,
	)

	forms = append(
		forms,
		buildThirdDeclensionForms(lemma, stem, "neuter", endings)...,
	)

	return forms
}

func buildThirdDeclensionForms(
	lemma models.Lemma,
	stem string,
	gender string,
	endings map[string]map[string]string,
) []models.Form {

	var forms []models.Form

	degree := "positive"

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
		for _, grammaticalCase := range cases {

			if (number == "singular" && grammaticalCase == "nominative") {
				forms = append(forms, models.Form{
					Form: lemma.Lemma,

					PartOfSpeech: "adjective",

					GrammaticalCase: StringPtr("nominative"),
					Number: "singular",
					Gender: &gender,

					Degree: &degree,
				})
			}

			forms = append(forms, models.Form{
				Form: stem + endings[number][grammaticalCase],

				PartOfSpeech: "adjective",

				GrammaticalCase: &grammaticalCase,
				Number: number,
				Gender: &gender,

				Degree: &degree,
			})
		}
	}

	return forms
}

// COMPARATIVE FORMS
func buildComparativeForms(stem string) []models.Form {

	var forms []models.Form

	// =====================================================
	// MASCULINE + FEMININE
	// =====================================================

	mfEndings := map[string]map[string]string{
		"singular": {
			"nominative": "ior",
			"genitive":   "iōris",
			"dative":     "iōrī",
			"accusative": "iōrem",
			"ablative":   "iōre",
			"vocative":   "ior",
		},
		"plural": {
			"nominative": "iōrēs",
			"genitive":   "iōrum",
			"dative":     "iōribus",
			"accusative": "iōrēs",
			"ablative":   "iōribus",
			"vocative":   "iōrēs",
		},
	}

	// masculine
	forms = append(
		forms,
		buildComparativeGenderForms(
			stem,
			"masculine",
			mfEndings,
		)...,
	)

	// feminine
	forms = append(
		forms,
		buildComparativeGenderForms(
			stem,
			"feminine",
			mfEndings,
		)...,
	)

	// =====================================================
	// NEUTER
	// =====================================================

	neuterEndings := map[string]map[string]string{
		"singular": {
			"nominative": "ius",
			"genitive":   "iōris",
			"dative":     "iōrī",
			"accusative": "ius",
			"ablative":   "iōre",
			"vocative":   "ius",
		},
		"plural": {
			"nominative": "iora",
			"genitive":   "iōrum",
			"dative":     "iōribus",
			"accusative": "iora",
			"ablative":   "iōribus",
			"vocative":   "iora",
		},
	}

	forms = append(
		forms,
		buildComparativeGenderForms(
			stem,
			"neuter",
			neuterEndings,
		)...,
	)

	return forms
}

func buildComparativeGenderForms(
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

	degree := "comparative"

	for _, number := range numbers {

		for _, grammaticalCase := range cases {

			form := stem + endings[number][grammaticalCase]

			forms = append(forms, models.Form{
				Form: form,

				PartOfSpeech: "adjective",

				GrammaticalCase: &grammaticalCase,
				Number:          number,
				Gender:          &gender,

				Degree: &degree,
			})
		}
	}

	return forms
}

func buildSuperlativeForms(stem string) []models.Form {

	superlativeStem := stem + "issim"

	var forms []models.Form

	forms = append(
		forms,
		buildSuperlativeGenderForms(
			superlativeStem,
			"masculine",
			map[string]map[string]string{
				"singular": {
					"nominative": "us",
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
			},
		)...,
	)

	forms = append(
		forms,
		buildSuperlativeGenderForms(
			superlativeStem,
			"feminine",
			map[string]map[string]string{
				"singular": {
					"nominative": "a",
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
			},
		)...,
	)

	forms = append(
		forms,
		buildSuperlativeGenderForms(
			superlativeStem,
			"neuter",
			map[string]map[string]string{
				"singular": {
					"nominative": "um",
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
			},
		)...,
	)

	return forms
}

func buildSuperlativeGenderForms(
	stem string,
	gender string,
	endings map[string]map[string]string,
) []models.Form {

	var forms []models.Form

	degree := "superlative"

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

				PartOfSpeech: "adjective",

				GrammaticalCase: &grammaticalCase,
				Number:          number,
				Gender:          &gender,

				Degree: &degree,
			})
		}
	}

	return forms
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

	degree := "positive"

	for _, number := range numbers {

		for _, grammaticalCase := range cases {

			forms = append(forms, models.Form{
				Form: stem + endings[number][grammaticalCase],

				PartOfSpeech: "adjective",

				GrammaticalCase: &grammaticalCase,

				Number: number,

				Gender: &gender,

				Degree: &degree,
			})
		}
	}

	return forms
}