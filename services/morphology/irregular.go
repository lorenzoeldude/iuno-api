package morphology

import "iuno-api/models"

//
// IRREGULAR VERB SYSTEM
//
// Latin has several highly irregular verbs
// that cannot be generated from normal
// conjugation rules.
//

func getIrregularForms(word models.Word) []models.Form {

	switch word.Lemma {

	case "esse":
		return generateEsse()

	case "posse":
		return generatePosse()

	case "velle":
		return generateVelle()

	case "ire":
		return generateIre()

	case "ferre":
		return generateFerre()
	}

	return []models.Form{}
}

// =====================================================
// ESSE
// =====================================================

func generateEsse() []models.Form {

	return []models.Form{

		// PRESENT

		{
			Form: "sum",
			Part: "verb",

			Number: "singular",

			Tense: "present",
			Mood:  "indicative",
			Voice: "active",

			Person: 1,
		},

		{
			Form: "es",
			Part: "verb",

			Number: "singular",

			Tense: "present",
			Mood:  "indicative",
			Voice: "active",

			Person: 2,
		},

		{
			Form: "est",
			Part: "verb",

			Number: "singular",

			Tense: "present",
			Mood:  "indicative",
			Voice: "active",

			Person: 3,
		},

		{
			Form: "sumus",
			Part: "verb",

			Number: "plural",

			Tense: "present",
			Mood:  "indicative",
			Voice: "active",

			Person: 1,
		},

		{
			Form: "estis",
			Part: "verb",

			Number: "plural",

			Tense: "present",
			Mood:  "indicative",
			Voice: "active",

			Person: 2,
		},

		{
			Form: "sunt",
			Part: "verb",

			Number: "plural",

			Tense: "present",
			Mood:  "indicative",
			Voice: "active",

			Person: 3,
		},

		// IMPERFECT

		{
			Form: "eram",
			Part: "verb",

			Number: "singular",

			Tense: "imperfect",
			Mood:  "indicative",
			Voice: "active",

			Person: 1,
		},

		{
			Form: "eras",
			Part: "verb",

			Number: "singular",

			Tense: "imperfect",
			Mood:  "indicative",
			Voice: "active",

			Person: 2,
		},

		{
			Form: "erat",
			Part: "verb",

			Number: "singular",

			Tense: "imperfect",
			Mood:  "indicative",
			Voice: "active",

			Person: 3,
		},

		{
			Form: "eramus",
			Part: "verb",

			Number: "plural",

			Tense: "imperfect",
			Mood:  "indicative",
			Voice: "active",

			Person: 1,
		},

		{
			Form: "eratis",
			Part: "verb",

			Number: "plural",

			Tense: "imperfect",
			Mood:  "indicative",
			Voice: "active",

			Person: 2,
		},

		{
			Form: "erant",
			Part: "verb",

			Number: "plural",

			Tense: "imperfect",
			Mood:  "indicative",
			Voice: "active",

			Person: 3,
		},

		// FUTURE

		{
			Form: "ero",
			Part: "verb",

			Number: "singular",

			Tense: "future",
			Mood:  "indicative",
			Voice: "active",

			Person: 1,
		},

		{
			Form: "eris",
			Part: "verb",

			Number: "singular",

			Tense: "future",
			Mood:  "indicative",
			Voice: "active",

			Person: 2,
		},

		{
			Form: "erit",
			Part: "verb",

			Number: "singular",

			Tense: "future",
			Mood:  "indicative",
			Voice: "active",

			Person: 3,
		},

		{
			Form: "erimus",
			Part: "verb",

			Number: "plural",

			Tense: "future",
			Mood:  "indicative",
			Voice: "active",

			Person: 1,
		},

		{
			Form: "eritis",
			Part: "verb",

			Number: "plural",

			Tense: "future",
			Mood:  "indicative",
			Voice: "active",

			Person: 2,
		},

		{
			Form: "erunt",
			Part: "verb",

			Number: "plural",

			Tense: "future",
			Mood:  "indicative",
			Voice: "active",

			Person: 3,
		},

		// INFINITIVES

		{
			Form: "esse",
			Part: "verb",

			Tense: "present",
			Voice: "active",

			NonFinite: "infinitive",
		},

		{
			Form: "fuisse",
			Part: "verb",

			Tense: "perfect",
			Voice: "active",

			NonFinite: "infinitive",
		},
	}
}

// =====================================================
// POSSE
// =====================================================

func generatePosse() []models.Form {

	return []models.Form{

		{
			Form: "possum",
			Part: "verb",

			Number: "singular",

			Tense: "present",
			Mood:  "indicative",
			Voice: "active",

			Person: 1,
		},

		{
			Form: "potes",
			Part: "verb",

			Number: "singular",

			Tense: "present",
			Mood:  "indicative",
			Voice: "active",

			Person: 2,
		},

		{
			Form: "potest",
			Part: "verb",

			Number: "singular",

			Tense: "present",
			Mood:  "indicative",
			Voice: "active",

			Person: 3,
		},

		{
			Form: "possumus",
			Part: "verb",

			Number: "plural",

			Tense: "present",
			Mood:  "indicative",
			Voice: "active",

			Person: 1,
		},

		{
			Form: "potestis",
			Part: "verb",

			Number: "plural",

			Tense: "present",
			Mood:  "indicative",
			Voice: "active",

			Person: 2,
		},

		{
			Form: "possunt",
			Part: "verb",

			Number: "plural",

			Tense: "present",
			Mood:  "indicative",
			Voice: "active",

			Person: 3,
		},

		{
			Form: "posse",
			Part: "verb",

			Tense: "present",
			Voice: "active",

			NonFinite: "infinitive",
		},
	}
}

// =====================================================
// VELLE
// =====================================================

func generateVelle() []models.Form {

	return []models.Form{

		{
			Form: "volo",
			Part: "verb",
			Number: "singular",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 1,
		},

		{
			Form: "vis",
			Part: "verb",
			Number: "singular",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 2,
		},

		{
			Form: "vult",
			Part: "verb",
			Number: "singular",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 3,
		},

		{
			Form: "volumus",
			Part: "verb",
			Number: "plural",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 1,
		},

		{
			Form: "vultis",
			Part: "verb",
			Number: "plural",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 2,
		},

		{
			Form: "volunt",
			Part: "verb",
			Number: "plural",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 3,
		},

		{
			Form: "velle",
			Part: "verb",
			Tense: "present",
			Voice: "active",
			NonFinite: "infinitive",
		},
	}
}

// =====================================================
// IRE
// =====================================================

func generateIre() []models.Form {

	return []models.Form{

		{
			Form: "eo",
			Part: "verb",
			Number: "singular",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 1,
		},

		{
			Form: "is",
			Part: "verb",
			Number: "singular",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 2,
		},

		{
			Form: "it",
			Part: "verb",
			Number: "singular",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 3,
		},

		{
			Form: "imus",
			Part: "verb",
			Number: "plural",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 1,
		},

		{
			Form: "itis",
			Part: "verb",
			Number: "plural",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 2,
		},

		{
			Form: "eunt",
			Part: "verb",
			Number: "plural",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 3,
		},

		{
			Form: "ire",
			Part: "verb",
			Tense: "present",
			Voice: "active",
			NonFinite: "infinitive",
		},
	}
}

// =====================================================
// FERRE
// =====================================================

func generateFerre() []models.Form {

	return []models.Form{

		{
			Form: "fero",
			Part: "verb",
			Number: "singular",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 1,
		},

		{
			Form: "fers",
			Part: "verb",
			Number: "singular",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 2,
		},

		{
			Form: "fert",
			Part: "verb",
			Number: "singular",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 3,
		},

		{
			Form: "ferimus",
			Part: "verb",
			Number: "plural",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 1,
		},

		{
			Form: "fertis",
			Part: "verb",
			Number: "plural",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 2,
		},

		{
			Form: "ferunt",
			Part: "verb",
			Number: "plural",
			Tense: "present",
			Mood: "indicative",
			Voice: "active",
			Person: 3,
		},

		{
			Form: "ferre",
			Part: "verb",
			Tense: "present",
			Voice: "active",
			NonFinite: "infinitive",
		},
	}
}