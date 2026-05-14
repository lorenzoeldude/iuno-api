package morphology

import "iuno-api/models"

//
// NON-FINITE VERB GENERATION
//
// This file generates:
//
// - infinitives
// - participles
// - gerunds
// - gerundives
// - supines
//
// These forms do NOT inflect by person.
//

func generateNonFiniteForms(
	word models.Word,
	stems VerbStems,
) []models.Form {

	var forms []models.Form

	// =====================================================
	// INFINITIVES
	// =====================================================

	forms = append(
		forms,
		generateInfinitives(word, stems)...,
	)

	// =====================================================
	// PARTICIPLES
	// =====================================================

	forms = append(
		forms,
		generateParticiples(word, stems)...,
	)

	// =====================================================
	// GERUND
	// =====================================================

	forms = append(
		forms,
		generateGerund(word, stems)...,
	)

	// =====================================================
	// GERUNDIVE
	// =====================================================

	forms = append(
		forms,
		generateGerundive(word, stems)...,
	)

	// =====================================================
	// SUPINE
	// =====================================================

	forms = append(
		forms,
		generateSupine(word, stems)...,
	)

	return forms
}

// =====================================================
// INFINITIVES
// =====================================================

func generateInfinitives(
	word models.Word,
	stems VerbStems,
) []models.Form {

	var forms []models.Form

	// PRESENT ACTIVE INFINITIVE
	//
	// amare
	// monere
	// regere
	// audire
	//

	forms = append(forms, models.Form{
		Form: word.Lemma,

		Part: "verb",

		Voice: "active",

		Tense: "present",

		NonFinite: "infinitive",
	})

	// PERFECT ACTIVE INFINITIVE
	//
	// amavisse
	//

	if stems.Perfect != "" {

		forms = append(forms, models.Form{
			Form: stems.Perfect + "isse",

			Part: "verb",

			Voice: "active",

			Tense: "perfect",

			NonFinite: "infinitive",
		})
	}

	// FUTURE ACTIVE INFINITIVE
	//
	// amaturus esse
	//

	if stems.Supine != "" {

		forms = append(forms, models.Form{
			Form: stems.Supine + "urus esse",

			Part: "verb",

			Voice: "active",

			Tense: "future",

			NonFinite: "infinitive",
		})
	}

	return forms
}

// =====================================================
// PARTICIPLES
// =====================================================

func generateParticiples(
	word models.Word,
	stems VerbStems,
) []models.Form {

	var forms []models.Form

	// PRESENT ACTIVE PARTICIPLE
	//
	// amans
	//

	if stems.Present != "" {

		forms = append(forms, models.Form{
			Form: stems.Present + "ns",

			Part: "verb",

			Tense: "present",
			Voice: "active",

			NonFinite: "participle",
		})
	}

	// FUTURE ACTIVE PARTICIPLE
	//
	// amaturus
	//

	if stems.Supine != "" {

		forms = append(forms, models.Form{
			Form: stems.Supine + "urus",

			Part: "verb",

			Tense: "future",
			Voice: "active",

			NonFinite: "participle",
		})
	}

	// PERFECT PASSIVE PARTICIPLE
	//
	// amatus
	//

	if stems.Supine != "" {

		forms = append(forms, models.Form{
			Form: stems.Supine + "us",

			Part: "verb",

			Tense: "perfect",
			Voice: "passive",

			NonFinite: "participle",
		})
	}

	return forms
}

// =====================================================
// GERUND
// =====================================================

func generateGerund(
	word models.Word,
	stems VerbStems,
) []models.Form {

	var forms []models.Form

	if stems.Present == "" {
		return forms
	}

	forms = append(forms,

		models.Form{
			Form: stems.Present + "ndi",

			Part: "verb",

			Case: "genitive",

			NonFinite: "gerund",
		},

		models.Form{
			Form: stems.Present + "ndo",

			Part: "verb",

			Case: "dative/ablative",

			NonFinite: "gerund",
		},

		models.Form{
			Form: stems.Present + "ndum",

			Part: "verb",

			Case: "accusative",

			NonFinite: "gerund",
		},
	)

	return forms
}

// =====================================================
// GERUNDIVE
// =====================================================

func generateGerundive(
	word models.Word,
	stems VerbStems,
) []models.Form {

	var forms []models.Form

	if stems.Present == "" {
		return forms
	}

	forms = append(forms, models.Form{
		Form: stems.Present + "ndus",

		Part: "verb",

		Voice: "passive",

		NonFinite: "gerundive",
	})

	return forms
}

// =====================================================
// SUPINE
// =====================================================

func generateSupine(
	word models.Word,
	stems VerbStems,
) []models.Form {

	var forms []models.Form

	if stems.Supine == "" {
		return forms
	}

	// accusative supine
	//
	// amatum
	//

	forms = append(forms, models.Form{
		Form: stems.Supine + "um",

		Part: "verb",

		Case: "accusative",

		NonFinite: "supine",
	})

	// ablative supine
	//
	// amatu
	//

	forms = append(forms, models.Form{
		Form: stems.Supine + "u",

		Part: "verb",

		Case: "ablative",

		NonFinite: "supine",
	})

	return forms
}