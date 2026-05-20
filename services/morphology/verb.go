package morphology

import (
	"iuno-api/models"
)

//
// VERB MORPHOLOGY ENGINE
//
// Supports:
// - 1st conjugation
// - Indicative
// - Subjunctive
// - Active voice
// - Present system + Perfect system
//

// type VerbForm struct {
// 	Form   string
// 	Person string
// 	Number string
// 	Tense  string
// 	Mood   string
// 	Voice  string
// }

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
	)...)

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
			Voice:  "active",
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

func removeVerbEnding(word string, ending string) string {
	if len(word) < len(ending) {
		return word
	}

	return word[:len(word)-len(ending)]
}