package morphology

import (
	// "log"
	"iuno-api/models"
	"strings"
)

func GenerateAdjective(lemma models.Lemma) []models.Form {

	if lemma.Declension == nil {
		return []models.Form{}
	}

	var forms []models.Form

	var stem string

	switch *lemma.Declension {

	case 12:
		stem = removeEnding(*lemma.Genitive, "ī")
		forms = append(forms,
            generateFirstSecondDeclensionAdjective(lemma, stem)...)

	case 31, 32, 33: 
		stem = removeEnding(*lemma.Genitive, "is")
		forms = append(forms,
            generateThirdDeclensionAdjective(lemma, stem)...)
	}

	forms = append(forms, buildComparativeForms(stem)...)

	forms = append(forms, buildSuperlativeForms(stem, lemma.Lemma)...)

	return forms
}

// =====================================================
// 1ST / 2ND DECLENSION ADJECTIVES
// =====================================================

func generateFirstSecondDeclensionAdjective(
	lemma models.Lemma,
	stem string,
) []models.Form {

	// stem := removeEnding(*lemma.Genitive, "ī")

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

	forms = append(forms, models.Form{
		LemmaID:      lemma.ID,
		PartOfSpeech: "adjective",
		Form: stem + "ē",
		Degree: StringPtr("positive"),
		FormType: StringPtr("adverb"),
	})

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
	stem string,
) []models.Form {

	endings := map[string]map[string]string{

		"singular": {
			"genitive":   "is",
			"dative":     "ī",

			// masculine/feminine only
			"accusative": "em",

			"ablative": "ī",
		},

		"plural": {
			// masculine/feminine only
			"nominative": "ēs",

			"genitive": "ium",
			"dative":   "ibus",

			// masculine/feminine only
			"accusative": "ēs",

			"ablative": "ibus",

			// masculine/feminine only
			"vocative": "ēs",
		},
	}

	var forms []models.Form

	forms = append(
		forms,
		buildThirdDeclensionForms(
			lemma,
			stem,
			"masculine",
			endings,
		)...,
	)

	forms = append(
		forms,
		buildThirdDeclensionForms(
			lemma,
			stem,
			"feminine",
			endings,
		)...,
	)

	forms = append(
		forms,
		buildThirdDeclensionForms(
			lemma,
			stem,
			"neuter",
			endings,
		)...,
	)

	var adverb string

	if strings.HasSuffix(lemma.Lemma, "ns") {
		adverb = stem + "er"
	} else {
		adverb = stem + "iter"
	}

	forms = append(forms, models.Form{
		LemmaID:      lemma.ID,
		PartOfSpeech: "adjective",
		Form: adverb,
		Degree: StringPtr("positive"),
		FormType: StringPtr("adverb"),
	})

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

			// =====================================
			// NOMINATIVE / VOCATIVE SINGULAR
			// =====================================

			if number == "singular" &&
				(grammaticalCase == "nominative" ||
					grammaticalCase == "vocative") {

				var form string

				switch *lemma.Declension {

				// one termination
				case 31:
					form = lemma.Lemma

				// two termination
				case 32:

					if gender == "neuter" {
						form = *lemma.Neuter
					} else {
						form = lemma.Lemma
					}

				// three termination
				case 33:

					switch gender {
					case "masculine":
						form = lemma.Lemma

					case "feminine":
						form = *lemma.Feminine

					case "neuter":
						form = *lemma.Neuter
					}
				}

				forms = append(forms, models.Form{
					Form: form,

					PartOfSpeech: "adjective",

					GrammaticalCase: &grammaticalCase,
					Number:          number,
					Gender:          &gender,

					Degree: &degree,

					FormType: StringPtr("adjective"),
				})

				continue
			}

			// =====================================
			// NEUTER SPECIAL FORMS
			// =====================================

			if gender == "neuter" {

				// neuter accusative singular
				if number == "singular" &&
					grammaticalCase == "accusative" {

					forms = append(forms, models.Form{
						Form: *lemma.Neuter,

						PartOfSpeech: "adjective",

						GrammaticalCase: &grammaticalCase,
						Number:          number,
						Gender:          &gender,

						Degree: &degree,

						FormType: StringPtr("adjective"),
					})

					continue
				}

				// neuter nominative/accusative/vocative plural
				if number == "plural" &&
					(grammaticalCase == "nominative" ||
						grammaticalCase == "accusative" ||
						grammaticalCase == "vocative") {

					forms = append(forms, models.Form{
						Form: stem + "ia",

						PartOfSpeech: "adjective",

						GrammaticalCase: &grammaticalCase,
						Number:          number,
						Gender:          &gender,

						Degree: &degree,

						FormType: StringPtr("adjective"),
					})

					continue
				}
			}

			// =====================================
			// REGULAR FORM
			// =====================================

			ending, ok := endings[number][grammaticalCase]
			if !ok {
				continue
			}

			forms = append(forms, models.Form{
				Form: stem + ending,

				PartOfSpeech: "adjective",

				GrammaticalCase: &grammaticalCase,
				Number:          number,
				Gender:          &gender,

				Degree: &degree,

				FormType: StringPtr("adjective"),
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
			"accusative": "iōra",
			"ablative":   "iōribus",
			"vocative":   "iōra",
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

	forms = append(forms, models.Form{
		Form: stem + "ius",

		PartOfSpeech: "adjective",

		Degree:   StringPtr("comparative"),
		FormType: StringPtr("adverb"),
	})

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

				FormType: StringPtr("adjective"),
			})
		}
	}

	return forms
}

func buildSuperlativeForms(stem string, lemma string) []models.Form {

	var superlativeStem string

	if strings.HasSuffix(lemma, "er") {
		superlativeStem = lemma + "rim"
	} else {
		superlativeStem = stem + "issim"
	}


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

	forms = append(forms, models.Form{
		Form: superlativeStem + "ē",

		PartOfSpeech: "adjective",

		Degree:   StringPtr("superlative"),
		FormType: StringPtr("adverb"),
	})

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

				FormType: StringPtr("adjective"),
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

				FormType: StringPtr("adjective"),
			})
		}
	}

	return forms
}