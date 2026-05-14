package morphology

import "iuno-api/models"

//
// FINITE VERB GENERATION
//
// This file generates:
//
// - indicative active
// - indicative passive
// - subjunctive active
// - subjunctive passive
//
// across:
//
// - person
// - number
// - tense
//
// This is the core of the morphology engine.
//

func generateFiniteForms(
	word models.Word,
	stems VerbStems,
) []models.Form {

	var forms []models.Form

	// =====================================================
	// ACTIVE INDICATIVE
	// =====================================================

	forms = append(
		forms,
		generatePresentActiveIndicative(word, stems)...,
	)

	forms = append(
		forms,
		generateImperfectActiveIndicative(word, stems)...,
	)

	forms = append(
		forms,
		generateFutureActiveIndicative(word, stems)...,
	)

	forms = append(
		forms,
		generatePerfectActiveIndicative(word, stems)...,
	)

	forms = append(
		forms,
		generatePluperfectActiveIndicative(word, stems)...,
	)

	forms = append(
		forms,
		generateFuturePerfectActiveIndicative(word, stems)...,
	)

	return forms
}

// =====================================================
// PRESENT ACTIVE INDICATIVE
// =====================================================

func generatePresentActiveIndicative(
	word models.Word,
	stems VerbStems,
) []models.Form {

	var forms []models.Form

	endings := PresentActiveEndings[word.Conjugation]

	for person := 1; person <= 3; person++ {

		// singular
		forms = append(forms, models.Form{
			Form: stems.Present + endings[personKey(person, "singular")],

			Part: "verb",

			Tense: "present",
			Mood:  "indicative",
			Voice: "active",

			Person: person,
			Number: "singular",
		})

		// plural
		forms = append(forms, models.Form{
			Form: stems.Present + endings[personKey(person, "plural")],

			Part: "verb",

			Tense: "present",
			Mood:  "indicative",
			Voice: "active",

			Person: person,
			Number: "plural",
		})
	}

	return forms
}

// =====================================================
// IMPERFECT ACTIVE INDICATIVE
// =====================================================

func generateImperfectActiveIndicative(
	word models.Word,
	stems VerbStems,
) []models.Form {

	var forms []models.Form

	endings := ImperfectActiveEndings[word.Conjugation]

	for person := 1; person <= 3; person++ {

		forms = append(forms, models.Form{
			Form: stems.Present + endings[personKey(person, "singular")],

			Part: "verb",

			Tense: "imperfect",
			Mood:  "indicative",
			Voice: "active",

			Person: person,
			Number: "singular",
		})

		forms = append(forms, models.Form{
			Form: stems.Present + endings[personKey(person, "plural")],

			Part: "verb",

			Tense: "imperfect",
			Mood:  "indicative",
			Voice: "active",

			Person: person,
			Number: "plural",
		})
	}

	return forms
}

// =====================================================
// FUTURE ACTIVE INDICATIVE
// =====================================================

func generateFutureActiveIndicative(
	word models.Word,
	stems VerbStems,
) []models.Form {

	var forms []models.Form

	endings := FutureActiveEndings[word.Conjugation]

	for person := 1; person <= 3; person++ {

		forms = append(forms, models.Form{
			Form: stems.Present + endings[personKey(person, "singular")],

			Part: "verb",

			Tense: "future",
			Mood:  "indicative",
			Voice: "active",

			Person: person,
			Number: "singular",
		})

		forms = append(forms, models.Form{
			Form: stems.Present + endings[personKey(person, "plural")],

			Part: "verb",

			Tense: "future",
			Mood:  "indicative",
			Voice: "active",

			Person: person,
			Number: "plural",
		})
	}

	return forms
}

// =====================================================
// PERFECT ACTIVE INDICATIVE
// =====================================================

func generatePerfectActiveIndicative(
	word models.Word,
	stems VerbStems,
) []models.Form {

	var forms []models.Form

	for person := 1; person <= 3; person++ {

		forms = append(forms, models.Form{
			Form: stems.Perfect + PerfectActiveEndings[personKey(person, "singular")],

			Part: "verb",

			Tense: "perfect",
			Mood:  "indicative",
			Voice: "active",

			Person: person,
			Number: "singular",
		})

		forms = append(forms, models.Form{
			Form: stems.Perfect + PerfectActiveEndings[personKey(person, "plural")],

			Part: "verb",

			Tense: "perfect",
			Mood:  "indicative",
			Voice: "active",

			Person: person,
			Number: "plural",
		})
	}

	return forms
}

// =====================================================
// PLUPERFECT ACTIVE INDICATIVE
// =====================================================

func generatePluperfectActiveIndicative(
	word models.Word,
	stems VerbStems,
) []models.Form {

	var forms []models.Form

	for person := 1; person <= 3; person++ {

		forms = append(forms, models.Form{
			Form: stems.Perfect + PluperfectActiveEndings[personKey(person, "singular")],

			Part: "verb",

			Tense: "pluperfect",
			Mood:  "indicative",
			Voice: "active",

			Person: person,
			Number: "singular",
		})

		forms = append(forms, models.Form{
			Form: stems.Perfect + PluperfectActiveEndings[personKey(person, "plural")],

			Part: "verb",

			Tense: "pluperfect",
			Mood:  "indicative",
			Voice: "active",

			Person: person,
			Number: "plural",
		})
	}

	return forms
}

// =====================================================
// FUTURE PERFECT ACTIVE INDICATIVE
// =====================================================

func generateFuturePerfectActiveIndicative(
	word models.Word,
	stems VerbStems,
) []models.Form {

	var forms []models.Form

	for person := 1; person <= 3; person++ {

		forms = append(forms, models.Form{
			Form: stems.Perfect + FuturePerfectActiveEndings[personKey(person, "singular")],

			Part: "verb",

			Tense: "future perfect",
			Mood:  "indicative",
			Voice: "active",

			Person: person,
			Number: "singular",
		})

		forms = append(forms, models.Form{
			Form: stems.Perfect + FuturePerfectActiveEndings[personKey(person, "plural")],

			Part: "verb",

			Tense: "future perfect",
			Mood:  "indicative",
			Voice: "active",

			Person: person,
			Number: "plural",
		})
	}

	return forms
}