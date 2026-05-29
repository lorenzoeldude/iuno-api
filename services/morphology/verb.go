package morphology

import (
	"iuno-api/models"
)

//
// VERB MORPHOLOGY ENGINE
//

func GenerateVerb(lemma models.Lemma) []models.Form {

	switch *lemma.Conjugation {
	case 1:
		return generateFirstConjugation(lemma)
	}

	return []models.Form{}
}

//
// =====================================================
// 1ST CONJUGATION
// =====================================================
//

func generateFirstConjugation(lemma models.Lemma) []models.Form {

	var forms []models.Form

	presentStem := removeVerbEnding(*lemma.Infinitive, "āre")
	perfectStem := removeVerbEnding(*lemma.Perfect, "ī")

	//
	// PRESENT ACTIVE INDICATIVE
	//

	forms = append(forms, buildVerbForms(
		lemma,
		presentStem,
		map[string]map[string]string{
			"singular": {
				"first":  "ō",
				"second": "ās",
				"third":  "at",
			},
			"plural": {
				"first":  "āmus",
				"second": "ātis",
				"third":  "ant",
			},
		},
		"present",
		"indicative",
		"active",
	)...)

	//
	// IMPERFECT ACTIVE INDICATIVE
	//

	forms = append(forms, buildVerbForms(
		lemma,
		presentStem,
		map[string]map[string]string{
			"singular": {
				"first":  "ābam",
				"second": "ābās",
				"third":  "ābat",
			},
			"plural": {
				"first":  "ābāmus",
				"second": "ābātis",
				"third":  "ābant",
			},
		},
		"imperfect",
		"indicative",
		"active",
	)...)

	//
	// FUTURE ACTIVE INDICATIVE
	//

	forms = append(forms, buildVerbForms(
		lemma,
		presentStem,
		map[string]map[string]string{
			"singular": {
				"first":  "ābō",
				"second": "ābis",
				"third":  "ābit",
			},
			"plural": {
				"first":  "ābimus",
				"second": "ābitis",
				"third":  "ābunt",
			},
		},
		"future",
		"indicative",
		"active",
	)...)

	//
	// PERFECT ACTIVE INDICATIVE
	//

	forms = append(forms, buildVerbForms(
		lemma,
		perfectStem,
		map[string]map[string]string{
			"singular": {
				"first":  "ī",
				"second": "istī",
				"third":  "it",
			},
			"plural": {
				"first":  "imus",
				"second": "istis",
				"third":  "ērunt",
			},
		},
		"perfect",
		"indicative",
		"active",
	)...)

	//
	// PLUPERFECT ACTIVE INDICATIVE
	//

	forms = append(forms, buildVerbForms(
		lemma,
		perfectStem,
		map[string]map[string]string{
			"singular": {
				"first":  "eram",
				"second": "erās",
				"third":  "erat",
			},
			"plural": {
				"first":  "erāmus",
				"second": "erātis",
				"third":  "erant",
			},
		},
		"pluperfect",
		"indicative",
		"active",
	)...)

	//
	// FUTURE PERFECT ACTIVE INDICATIVE
	//

	forms = append(forms, buildVerbForms(
		lemma,
		perfectStem,
		map[string]map[string]string{
			"singular": {
				"first":  "erō",
				"second": "eris",
				"third":  "erit",
			},
			"plural": {
				"first":  "erimus",
				"second": "eritis",
				"third":  "erint",
			},
		},
		"future perfect",
		"indicative",
		"active",
	)...)

	//
	// PRESENT ACTIVE SUBJUNCTIVE
	//

	forms = append(forms, buildVerbForms(
		lemma,
		presentStem,
		map[string]map[string]string{
			"singular": {
				"first":  "em",
				"second": "ēs",
				"third":  "et",
			},
			"plural": {
				"first":  "ēmus",
				"second": "ētis",
				"third":  "ent",
			},
		},
		"present",
		"subjunctive",
		"active",
	)...)

	//
	// IMPERFECT ACTIVE SUBJUNCTIVE
	//

	imperfectSubjStem := removeVerbEnding(*lemma.Infinitive, "re")

	forms = append(forms, buildVerbForms(
		lemma,
		imperfectSubjStem,
		map[string]map[string]string{
			"singular": {
				"first":  "m",
				"second": "s",
				"third":  "t",
			},
			"plural": {
				"first":  "mus",
				"second": "tis",
				"third":  "nt",
			},
		},
		"imperfect",
		"subjunctive",
		"active",
	)...)

	//
	// PERFECT ACTIVE SUBJUNCTIVE
	//

	forms = append(forms, buildVerbForms(
		lemma,
		perfectStem,
		map[string]map[string]string{
			"singular": {
				"first":  "erim",
				"second": "erīs",
				"third":  "erit",
			},
			"plural": {
				"first":  "erīmus",
				"second": "erītis",
				"third":  "erint",
			},
		},
		"perfect",
		"subjunctive",
		"active",
	)...)

	//
	// PLUPERFECT ACTIVE SUBJUNCTIVE
	//

	pluperfectSubjStem := perfectStem + "isse"

	forms = append(forms, buildVerbForms(
		lemma,
		pluperfectSubjStem,
		map[string]map[string]string{
			"singular": {
				"first":  "m",
				"second": "s",
				"third":  "t",
			},
			"plural": {
				"first":  "mus",
				"second": "tis",
				"third":  "nt",
			},
		},
		"pluperfect",
		"subjunctive",
		"active",
	)...)

	//
	// PRESENT ACTIVE IMPERATIVE
	//

	forms = append(forms,
		models.Form{
			Form:   presentStem + "ā",
			Part:   "verb",
			Person: 2,
			Number: "singular",
			Tense:  "present",
			Mood:   "imperative",
			Voice:  "active",
		},
		models.Form{
			Form:   presentStem + "āte",
			Part:   "verb",
			Person: 2,
			Number: "plural",
			Tense:  "present",
			Mood:   "imperative",
			Voice:  "active",
		},
	)

	//
	// PRESENT PASSIVE INDICATIVE
	//
	forms = append(forms, buildVerbForms(
		lemma,
		presentStem,
		map[string]map[string]string{
			"singular": {
				"first":  "or",
				"second": "āris",
				"third":  "ātur",
			},
			"plural": {
				"first":  "āmur",
				"second": "āminī",
				"third":  "antur",
			},
		},
		"present",
		"indicative",
		"passive",
	)...)

	//
	// IMPERFECT PASSIVE INDICATIVE
	//

	forms = append(forms, buildVerbForms(
		lemma,
		presentStem,
		map[string]map[string]string{
			"singular": {
				"first":  "ābar",
				"second": "ābāris",
				"third":  "ābātur",
			},
			"plural": {
				"first":  "ābāmur",
				"second": "ābāminī",
				"third":  "ābantur",
			},
		},
		"imperfect",
		"indicative",
		"passive",
	)...)

	//
	// FUTURE PASSIVE INDICATIVE
	//

	forms = append(forms, buildVerbForms(
		lemma,
		presentStem,
		map[string]map[string]string{
			"singular": {
				"first":  "ābor",
				"second": "āberis",
				"third":  "ābitur",
			},
			"plural": {
				"first":  "ābimur",
				"second": "ābiminī",
				"third":  "ābuntur",
			},
		},
		"future",
		"indicative",
		"passive",
	)...)

	/// PPP stem
	ppp := pppStem(*lemma.Supine)

	//
	// PERFECT PASSIVE INDICATIVE
	//

	forms = append(forms,

		models.Form{
			Form:   ppp + "us sum",
			Part:   "verb",
			Person: 1,
			Number: "singular",
			Tense:  "perfect",
			Mood:   "indicative",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "us es",
			Part:   "verb",
			Person: 2,
			Number: "singular",
			Tense:  "perfect",
			Mood:   "indicative",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "us est",
			Part:   "verb",
			Person: 3,
			Number: "singular",
			Tense:  "perfect",
			Mood:   "indicative",
			Voice:  "passive",
		},

		models.Form{
			Form:   ppp + "ī sumus",
			Part:   "verb",
			Person: 1,
			Number: "plural",
			Tense:  "perfect",
			Mood:   "indicative",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "ī estis",
			Part:   "verb",
			Person: 2,
			Number: "plural",
			Tense:  "perfect",
			Mood:   "indicative",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "ī sunt",
			Part:   "verb",
			Person: 3,
			Number: "plural",
			Tense:  "perfect",
			Mood:   "indicative",
			Voice:  "passive",
		},
	)

	//
	// PLUPERFECT PASSIVE INDICATIVE
	//

	forms = append(forms,

		models.Form{
			Form:   ppp + "us eram",
			Part:   "verb",
			Person: 1,
			Number: "singular",
			Tense:  "pluperfect",
			Mood:   "indicative",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "us erās",
			Part:   "verb",
			Person: 2,
			Number: "singular",
			Tense:  "pluperfect",
			Mood:   "indicative",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "us erat",
			Part:   "verb",
			Person: 3,
			Number: "singular",
			Tense:  "pluperfect",
			Mood:   "indicative",
			Voice:  "passive",
		},

		models.Form{
			Form:   ppp + "ī erāmus",
			Part:   "verb",
			Person: 1,
			Number: "plural",
			Tense:  "pluperfect",
			Mood:   "indicative",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "ī erātis",
			Part:   "verb",
			Person: 2,
			Number: "plural",
			Tense:  "pluperfect",
			Mood:   "indicative",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "ī erant",
			Part:   "verb",
			Person: 3,
			Number: "plural",
			Tense:  "pluperfect",
			Mood:   "indicative",
			Voice:  "passive",
		},
	)

	//
	// FUTURE PERFECT PASSIVE INDICATIVE
	//

	forms = append(forms,

		models.Form{
			Form:   ppp + "us erō",
			Part:   "verb",
			Person: 1,
			Number: "singular",
			Tense:  "future perfect",
			Mood:   "indicative",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "us eris",
			Part:   "verb",
			Person: 2,
			Number: "singular",
			Tense:  "future perfect",
			Mood:   "indicative",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "us erit",
			Part:   "verb",
			Person: 3,
			Number: "singular",
			Tense:  "future perfect",
			Mood:   "indicative",
			Voice:  "passive",
		},

		models.Form{
			Form:   ppp + "ī erimus",
			Part:   "verb",
			Person: 1,
			Number: "plural",
			Tense:  "future perfect",
			Mood:   "indicative",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "ī eritis",
			Part:   "verb",
			Person: 2,
			Number: "plural",
			Tense:  "future perfect",
			Mood:   "indicative",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "ī erunt",
			Part:   "verb",
			Person: 3,
			Number: "plural",
			Tense:  "future perfect",
			Mood:   "indicative",
			Voice:  "passive",
		},
	)

	//
	// PRESENT PASSIVE SUBJUNCTIVE
	//

	forms = append(forms, buildVerbForms(
		lemma,
		presentStem,
		map[string]map[string]string{
			"singular": {
				"first":  "er",
				"second": "ēris",
				"third":  "ētur",
			},
			"plural": {
				"first":  "ēmur",
				"second": "ēminī",
				"third":  "entur",
			},
		},
		"present",
		"subjunctive",
		"passive",
	)...)

	//
	// IMPERFECT PASSIVE SUBJUNCTIVE
	//
	impSubPasStem := removeVerbEnding(*lemma.Infinitive, "re")

	forms = append(forms, buildVerbForms(
		lemma,
		impSubPasStem,
		map[string]map[string]string{
			"singular": {
				"first":  "rer",
				"second": "rēris",
				"third":  "rētur",
			},
			"plural": {
				"first":  "rēmur",
				"second": "rēminī",
				"third":  "rentur",
			},
		},
		"imperfect",
		"subjunctive",
		"passive",
	)...)

	//
	// PERFECT PASSIVE SUBJUNCTIVE
	//

	forms = append(forms,

		models.Form{
			Form:   ppp + "us sim",
			Part:   "verb",
			Person: 1,
			Number: "singular",
			Tense:  "perfect",
			Mood:   "subjunctive",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "us sīs",
			Part:   "verb",
			Person: 2,
			Number: "singular",
			Tense:  "perfect",
			Mood:   "subjunctive",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "us sit",
			Part:   "verb",
			Person: 3,
			Number: "singular",
			Tense:  "perfect",
			Mood:   "subjunctive",
			Voice:  "passive",
		},

		models.Form{
			Form:   ppp + "ī sīmus",
			Part:   "verb",
			Person: 1,
			Number: "plural",
			Tense:  "perfect",
			Mood:   "subjunctive",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "ī sītis",
			Part:   "verb",
			Person: 2,
			Number: "plural",
			Tense:  "perfect",
			Mood:   "subjunctive",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "ī sint",
			Part:   "verb",
			Person: 3,
			Number: "plural",
			Tense:  "perfect",
			Mood:   "subjunctive",
			Voice:  "passive",
		},
	)

	//
	// PLUPERFECT PASSIVE SUBJUNCTIVE
	//

	forms = append(forms,

		models.Form{
			Form:   ppp + "us essem",
			Part:   "verb",
			Person: 1,
			Number: "singular",
			Tense:  "pluperfect",
			Mood:   "subjunctive",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "us essēs",
			Part:   "verb",
			Person: 2,
			Number: "singular",
			Tense:  "pluperfect",
			Mood:   "subjunctive",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "us esset",
			Part:   "verb",
			Person: 3,
			Number: "singular",
			Tense:  "pluperfect",
			Mood:   "subjunctive",
			Voice:  "passive",
		},

		models.Form{
			Form:   ppp + "ī essēmus",
			Part:   "verb",
			Person: 1,
			Number: "plural",
			Tense:  "pluperfect",
			Mood:   "subjunctive",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "ī essētis",
			Part:   "verb",
			Person: 2,
			Number: "plural",
			Tense:  "pluperfect",
			Mood:   "subjunctive",
			Voice:  "passive",
		},
		models.Form{
			Form:   ppp + "ī essent",
			Part:   "verb",
			Person: 3,
			Number: "plural",
			Tense:  "pluperfect",
			Mood:   "subjunctive",
			Voice:  "passive",
		},
	)

	//
	// PRESENT PASSIVE IMPERATIVE
	//

	forms = append(forms,
		models.Form{
			Form:   presentStem + "āre",
			Part:   "verb",
			Person: 2,
			Number: "singular",
			Tense:  "present",
			Mood:   "imperative",
			Voice:  "passive",
		},
		models.Form{
			Form:   presentStem + "āminī",
			Part:   "verb",
			Person: 2,
			Number: "plural",
			Tense:  "present",
			Mood:   "imperative",
			Voice:  "passive",
		},
	)

	return forms
}

//
// =====================================================
// BUILD VERB FORMS
// =====================================================
//

func buildVerbForms(
	lemma models.Lemma,
	stem string,
	endings map[string]map[string]string,
	tense string,
	mood string,
	voice string,
) []models.Form {

	var forms []models.Form

	numbers := []string{
		"singular",
		"plural",
	}

	// numbers := []int{1, 2}

	persons := []int{1, 2, 3}

	for _, number := range numbers {
	for _, person := range persons {

		var personKey string

		switch person {
		case 1:
			personKey = "first"
		case 2:
			personKey = "second"
		case 3:
			personKey = "third"
		}

		form := stem + endings[number][personKey]

		forms = append(forms, models.Form{
			Form:   form,
			Person: person,
			Number: number,
			Tense:  tense,
			Mood:   mood,
			Voice:  voice,
			Part: "verb",
		})
	}
}

	return forms
}

func buildPerfectPassiveForms(
	ppp string,
	sumForms map[string]map[string]string,
	tense string,
	mood string,
) []models.Form {

	var forms []models.Form

	numbers := []string{"singular", "plural"}
	persons := []int{1,2,3}

	for _, number := range numbers {
		for _, person := range persons {

			var key string

			switch person {
			case 1:
				key = "first"
			case 2:
				key = "second"
			case 3:
				key = "third"
			}

			form := ppp + "us " + sumForms[number][key]

			forms = append(forms, models.Form{
				Form: form,
				Person: person,
				Number: number,
				Tense: tense,
				Mood: mood,
				Voice: "passive",
				Part: "verb",
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

func pppStem(supine string) string {
	return removeEnding(supine, "um") 
}

func removeVerbEnding(word string, ending string) string {
	if len(word) < len(ending) {
		return word
	}

	return word[:len(word)-len(ending)]
}
