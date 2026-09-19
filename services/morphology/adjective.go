package morphology

import (
	"log"
	"strings"

	"iuno-api/models"
)

func GenerateAdjective(lemma models.Lemma) []models.Form {

	if lemma.Declension == nil {
		log.Printf("Skipping adjective %s: missing declension", lemma.Lemma)
		return []models.Form{}
	}

	if lemma.Genitive == nil {
		log.Printf("Skipping adjective %s: missing genitive", lemma.Lemma)
		return []models.Form{}
	}

	var forms []models.Form

	var stem string

	switch *lemma.Declension {

	case 12:
		isPronominal := strings.HasSuffix(*lemma.Genitive, "īus")

		if isPronominal {
			stem = removeEnding(*lemma.Genitive, "īus")
		} else {
			stem = removeEnding(*lemma.Genitive, "ī")
		}

		forms = append(
			forms,
			generateFirstSecondDeclensionAdjective(
				lemma,
				stem,
				isPronominal,
			)...,
		)

	case 31, 32, 33:
		stem = removeEnding(*lemma.Genitive, "is")

		forms = append(
			forms,
			generateThirdDeclensionAdjective(lemma, stem)...,
		)
	}

	// =====================================================
	// COMPARATIVE / SUPERLATIVE
	// =====================================================

	if lemma.Comparable != nil && *lemma.Comparable {

		if lemma.Comparative != nil {
			forms = append(
				forms,
				buildComparativeForms(*lemma.Comparative)...,
			)
		}

		if lemma.Superlative != nil {
			forms = append(
				forms,
				buildSuperlativeForms(*lemma.Superlative)...,
			)
		}
	}

	return forms
}

// =====================================================
// 1ST / 2ND DECLENSION ADJECTIVES
// =====================================================

func generateFirstSecondDeclensionAdjective(
	lemma models.Lemma,
	stem string,
	isPronominal bool,
) []models.Form {

	var forms []models.Form

	if isPronominal {
		forms = append(
			forms,
			generatePronominalAdjectiveForms(lemma, stem)...,
		)
	} else {
		forms = append(
			forms,
			buildMasculineAdjectiveForms(
				lemma.Lemma,
				stem,
			)...,
		)

		forms = append(
			forms,
			buildFeminineAdjectiveForms(
				lemma,
				stem,
			)...,
		)

		forms = append(
			forms,
			buildNeuterAdjectiveForms(
				lemma,
				stem,
			)...,
		)

		forms = append(forms, models.Form{
			LemmaID:      lemma.ID,
			PartOfSpeech: "adjective",
			Form:         stem + "ē",
			Degree:       StringPtr("positive"),
			FormType:     StringPtr("adverb"),
		})
	}

	return forms
}

func generatePronominalAdjectiveForms(
	lemma models.Lemma,
	stem string,
) []models.Form {

	var forms []models.Form

	masculineEndings := map[string]map[string]string{
		"singular": {
			"genitive":   "īus",
			"dative":     "ī",
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

	feminineEndings := map[string]map[string]string{
		"singular": {
			"nominative": "a",
			"genitive":   "īus",
			"dative":     "ī",
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

	neuterEndings := map[string]map[string]string{
		"singular": {
			"nominative": "um",
			"genitive":   "īus",
			"dative":     "ī",
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

	// Masculine nominative singular is stored as the lemma itself.
	forms = append(forms, models.Form{
		LemmaID:         lemma.ID,
		Form:            lemma.Lemma,
		PartOfSpeech:    "adjective",
		GrammaticalCase: StringPtr("nominative"),
		Number:          "singular",
		Gender:          StringPtr("masculine"),
		Degree:          StringPtr("positive"),
		FormType:        StringPtr("adjective"),
	})

	forms = append(
		forms,
		buildAdjectiveForms(
			stem,
			"masculine",
			masculineEndings,
		)...,
	)

	forms = append(
		forms,
		buildAdjectiveForms(
			stem,
			"feminine",
			feminineEndings,
		)...,
	)

	// Some pronominal adjectives have a special neuter singular.
	// Example: alius → aliud
	if lemma.Neuter != nil {
		neuterEndings["singular"]["nominative"] = ""
		neuterEndings["singular"]["accusative"] = ""
		neuterEndings["singular"]["vocative"] = ""
	}

	neuterForms := buildAdjectiveForms(
		stem,
		"neuter",
		neuterEndings,
	)

	if lemma.Neuter != nil {
		for i := range neuterForms {
			if neuterForms[i].Number == "singular" &&
				neuterForms[i].GrammaticalCase != nil &&
				(*neuterForms[i].GrammaticalCase == "nominative" ||
					*neuterForms[i].GrammaticalCase == "accusative" ||
					*neuterForms[i].GrammaticalCase == "vocative") {

				neuterForms[i].Form = *lemma.Neuter
			}
		}
	}

	forms = append(forms, neuterForms...)

	forms = append(forms, models.Form{
		LemmaID:      lemma.ID,
		PartOfSpeech: "adjective",
		Form:         stem + "ē",
		Degree:       StringPtr("positive"),
		FormType:     StringPtr("adverb"),
	})

	return forms
}

func buildMasculineAdjectiveForms(
	nominative string,
	stem string,
) []models.Form {

	endings := map[string]map[string]string{

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

	var forms []models.Form

	forms = append(forms, models.Form{
		Form:            nominative,
		PartOfSpeech:    "adjective",
		GrammaticalCase: StringPtr("nominative"),
		Number:          "singular",
		Gender:          StringPtr("masculine"),
		Degree:          StringPtr("positive"),
		FormType:        StringPtr("adjective"),
	})

	forms = append(
		forms,
		buildAdjectiveForms(
			stem,
			"masculine",
			endings,
		)...,
	)

	return forms
}

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

	forms := buildAdjectiveForms(
		stem,
		"neuter",
		endings,
	)

	return forms
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
			"accusative": "em",
			"ablative":   "ī",
		},

		"plural": {
			"nominative": "ēs",
			"genitive":   "ium",
			"dative":     "ibus",
			"accusative": "ēs",
			"ablative":   "ibus",
			"vocative":   "ēs",
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
		Form:         adverb,
		Degree:       StringPtr("positive"),
		FormType:     StringPtr("adverb"),
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

			// =====================================
			// NOMINATIVE / VOCATIVE SINGULAR
			// =====================================

			if number == "singular" &&
				(grammaticalCase == "nominative" ||
					grammaticalCase == "vocative") {

				var form string

				switch *lemma.Declension {

				// One termination.
				case 31:
					form = lemma.Lemma

				// Two terminations.
				case 32:

					if gender == "neuter" {
						if lemma.Neuter == nil {
							log.Printf(
								"Skipping adjective %s: missing neuter",
								lemma.Lemma,
							)
							continue
						}

						form = *lemma.Neuter
					} else {
						form = lemma.Lemma
					}

				// Three terminations.
				case 33:

					switch gender {

					case "masculine":
						form = lemma.Lemma

					case "feminine":
						if lemma.Feminine == nil {
							log.Printf(
								"Skipping adjective %s: missing feminine",
								lemma.Lemma,
							)
							continue
						}

						form = *lemma.Feminine

					case "neuter":
						if lemma.Neuter == nil {
							log.Printf(
								"Skipping adjective %s: missing neuter",
								lemma.Lemma,
							)
							continue
						}

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

				// Neuter accusative singular.
				if number == "singular" &&
					grammaticalCase == "accusative" {

					if lemma.Neuter == nil {
						log.Printf(
							"Skipping adjective %s: missing neuter",
							lemma.Lemma,
						)
						continue
					}

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

				// Neuter nominative/accusative/vocative plural.
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

// =====================================================
// COMPARATIVE FORMS
// =====================================================

func buildComparativeForms(
	comparative string,
) []models.Form {

	var forms []models.Form

	// Comparative masculine/feminine nominative singular:
	//
	// fortior
	// melior
	// maior
	//
	// The oblique stem is formed by replacing -ior
	// with -iōr.
	stem := removeEnding(comparative, "ior") + "iōr"

	mfEndings := map[string]map[string]string{
		"singular": {
			"nominative": "ior",
			"genitive":   "is",
			"dative":     "ī",
			"accusative": "em",
			"ablative":   "e",
			"vocative":   "ior",
		},
		"plural": {
			"nominative": "ēs",
			"genitive":   "um",
			"dative":     "ibus",
			"accusative": "ēs",
			"ablative":   "ibus",
			"vocative":   "ēs",
		},
	}

	// Masculine.
	forms = append(
		forms,
		buildComparativeGenderForms(
			comparative,
			stem,
			"masculine",
			mfEndings,
		)...,
	)

	// Feminine.
	forms = append(
		forms,
		buildComparativeGenderForms(
			comparative,
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
			"genitive": "is",
			"dative":   "ī",
			"ablative": "e",
		},
		"plural": {
			"genitive": "um",
			"dative":   "ibus",
			"ablative": "ibus",
		},
	}

	forms = append(
		forms,
		buildComparativeNeuterForms(
			comparative,
			stem,
			neuterEndings,
		)...,
	)

	// Comparative adverb.
	forms = append(forms, models.Form{
		Form:         stem + "ius",
		PartOfSpeech: "adjective",
		Degree:       StringPtr("comparative"),
		FormType:     StringPtr("adverb"),
	})

	return forms
}

func buildComparativeGenderForms(
	nominative string,
	stem string,
	gender string,
	endings map[string]map[string]string,
) []models.Form {

	var forms []models.Form

	degree := "comparative"

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

			var form string

			// Masculine/feminine nominative and vocative
			// singular use the stored comparative.
			if number == "singular" &&
				(grammaticalCase == "nominative" ||
					grammaticalCase == "vocative") {

				form = nominative
			} else {
				form = stem + endings[number][grammaticalCase]
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
		}
	}

	return forms
}

func buildComparativeNeuterForms(
	nominative string,
	stem string,
	endings map[string]map[string]string,
) []models.Form {

	var forms []models.Form

	degree := "comparative"

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

			var form string

			switch {

			// Neuter singular nominative/accusative/vocative.
			case number == "singular" &&
				(grammaticalCase == "nominative" ||
					grammaticalCase == "accusative" ||
					grammaticalCase == "vocative"):

				form = removeEnding(nominative, "ior") + "ius"

			// Neuter plural nominative/accusative/vocative.
			case number == "plural" &&
				(grammaticalCase == "nominative" ||
					grammaticalCase == "accusative" ||
					grammaticalCase == "vocative"):

				form = stem + "a"

			default:
				form = stem + endings[number][grammaticalCase]
			}

			forms = append(forms, models.Form{
				Form: form,

				PartOfSpeech: "adjective",

				GrammaticalCase: &grammaticalCase,
				Number:          number,
				Gender:          StringPtr("neuter"),

				Degree: &degree,

				FormType: StringPtr("adjective"),
			})
		}
	}

	return forms
}

// =====================================================
// SUPERLATIVE FORMS
// =====================================================

func buildSuperlativeForms(
	superlative string,
) []models.Form {

	// The stored superlative is the masculine nominative
	// singular, e.g.:
	//
	// optimus
	// fortissimus
	// pulcherrimus
	//
	// Remove -us to get the declension stem.
	stem := removeEnding(superlative, "us")

	var forms []models.Form

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

	forms = append(
		forms,
		buildSuperlativeGenderForms(
			superlative,
			stem,
			"masculine",
			endings,
		)...,
	)

	forms = append(
		forms,
		buildSuperlativeGenderForms(
			superlative,
			stem,
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
			superlative,
			stem,
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

	// Superlative adverb.
	forms = append(forms, models.Form{
		Form:         stem + "ē",
		PartOfSpeech: "adjective",
		Degree:       StringPtr("superlative"),
		FormType:     StringPtr("adverb"),
	})

	return forms
}

func buildSuperlativeGenderForms(
	nominative string,
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

			var form string

			// Masculine nominative/vocative singular.
			if gender == "masculine" &&
				number == "singular" &&
				(grammaticalCase == "nominative" ||
					grammaticalCase == "vocative") {

				form = nominative

			} else {
				form = stem + endings[number][grammaticalCase]
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
