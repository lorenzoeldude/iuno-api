package morphology

import (
	"iuno-api/models"
)

func infinitiveForm(
    form string,
    tense string,
    voice string,
) models.Form {
    return models.Form{
        Form: form,
        PartOfSpeech: "verb",
        Tense: StringPtr(tense),
        Mood: StringPtr("infinitive"),
        Voice: StringPtr(voice),
    }
}

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

	/// PPP stem
	ppp := pppStem(*lemma.Supine)

	// Active Infinitives
	forms = append(forms,
		infinitiveForm(presentStem+"āre", "present", "active"),
		infinitiveForm(perfectStem+"isse", "perfect", "active"),
		infinitiveForm(ppp+"ūrus esse", "future", "active"),
	)

	// Passive Infinitives
	forms = append(forms,
		infinitiveForm(presentStem+"ārī", "present", "passive"),
		infinitiveForm(ppp+"um esse", "perfect", "passive"),
		infinitiveForm(ppp+"um īrī", "future", "passive"),
	)

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
			PartOfSpeech:   "verb",
			Person: IntPtr(2),
			Number: "singular",
			Tense:  StringPtr("present"),
			Mood:   StringPtr("imperative"),
			Voice:  StringPtr("active"),
		},
		models.Form{
			Form:   presentStem + "āte",
			PartOfSpeech:   "verb",
			Person: IntPtr(2),
			Number: "plural",
			Tense:  StringPtr("present"),
			Mood:   StringPtr("imperative"),
			Voice:  StringPtr("active"),
		},
	)

	//
	// FUTURE ACTIVE IMPERATIVE
	//

	forms = append(forms,
		models.Form{
			Form: presentStem + "ātō",
			PartOfSpeech: "verb",
			Person: IntPtr(2),
			Number: "singular",
			Tense: StringPtr("future"),
			Mood: StringPtr("imperative"),
			Voice: StringPtr("active"),
		},
		models.Form{
			Form: presentStem + "ātōte",
			PartOfSpeech: "verb",
			Person: IntPtr(2),
			Number: "plural",
			Tense: StringPtr("future"),
			Mood: StringPtr("imperative"),
			Voice: StringPtr("active"),
		},
		models.Form{
			Form: presentStem + "ātō",
			PartOfSpeech: "verb",
			Person: IntPtr(3),
			Number: "singular",
			Tense: StringPtr("future"),
			Mood: StringPtr("imperative"),
			Voice: StringPtr("active"),
		},
		models.Form{
			Form: presentStem + "antō",
			PartOfSpeech: "verb",
			Person: IntPtr(3),
			Number: "plural",
			Tense: StringPtr("future"),
			Mood: StringPtr("imperative"),
			Voice: StringPtr("active"),
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

	//
	// PERFECT PASSIVE INDICATIVE
	//

	forms = append(forms, buildPerfectPassiveForms(
		ppp,
		map[string]map[string]string{
			"singular": {
				"first":  "sum",
				"second": "es",
				"third":  "est",
			},
			"plural": {
				"first":  "sumus",
				"second": "estis",
				"third":  "sunt",
			},
		},
		"perfect",
		"indicative",
	)...)

	//
	// PLUPERFECT PASSIVE INDICATIVE
	//
	forms = append(forms, buildPerfectPassiveForms(
		ppp,
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
	// FUTURE PERFECT PASSIVE INDICATIVE
	//
	forms = append(forms, buildPerfectPassiveForms(
		ppp,
		map[string]map[string]string{
			"singular": {
				"first":  "erō",
				"second": "eris",
				"third":  "erit",
			},
			"plural": {
				"first":  "erimus",
				"second": "eritis",
				"third":  "erunt",
			},
		},
		"future perfect",
		"indicative",
	)...)

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
	forms = append(forms, buildPerfectPassiveForms(
		ppp,
		map[string]map[string]string{
			"singular": {
				"first":  "sim",
				"second": "sīs",
				"third":  "sit",
			},
			"plural": {
				"first":  "sīmus",
				"second": "sītis",
				"third":  "sint",
			},
		},
		"perfect",
		"subjunctive",
	)...)

	//
	// PLUPERFECT PASSIVE SUBJUNCTIVE
	//
	forms = append(forms, buildPerfectPassiveForms(
		ppp,
		map[string]map[string]string{
			"singular": {
				"first":  "essem",
				"second": "essēs",
				"third":  "esset",
			},
			"plural": {
				"first":  "essēmus",
				"second": "essētis",
				"third":  "essent",
			},
		},
		"pluperfect",
		"subjunctive",
	)...)

	//
	// PRESENT PASSIVE IMPERATIVE
	//

	forms = append(forms,
		models.Form{
			Form:   presentStem + "āre",
			PartOfSpeech:   "verb",
			Person: IntPtr(2),
			Number: "singular",
			Tense:  StringPtr("present"),
			Mood:   StringPtr("imperative"),
			Voice:  StringPtr("passive"),
		},
		models.Form{
			Form:   presentStem + "āminī",
			PartOfSpeech:   "verb",
			Person: IntPtr(2),
			Number: "plural",
			Tense:  StringPtr("present"),
			Mood:   StringPtr("imperative"),
			Voice:  StringPtr("passive"),
		},
	)

	//
	// FUTURE PASSIVE IMPERATIVE
	//

	forms = append(forms,
		models.Form{
			Form: presentStem + "ātor",
			PartOfSpeech: "verb",
			Person: IntPtr(2),
			Number: "singular",
			Tense: StringPtr("future"),
			Mood: StringPtr("imperative"),
			Voice: StringPtr("passive"),
		},
		models.Form{
			Form: presentStem + "āminī",
			PartOfSpeech: "verb",
			Person: IntPtr(2),
			Number: "plural",
			Tense: StringPtr("future"),
			Mood: StringPtr("imperative"),
			Voice: StringPtr("passive"),
		},
		models.Form{
			Form: presentStem + "ātor",
			PartOfSpeech: "verb",
			Person: IntPtr(3),
			Number: "singular",
			Tense: StringPtr("future"),
			Mood: StringPtr("imperative"),
			Voice: StringPtr("passive"),
		},
		models.Form{
			Form: presentStem + "antor",
			PartOfSpeech: "verb",
			Person: IntPtr(3),
			Number: "plural",
			Tense: StringPtr("future"),
			Mood: StringPtr("imperative"),
			Voice: StringPtr("passive"),
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

			forms = append(forms,
				verbForm(
					form,
					person,
					number,
					tense,
					mood,
					voice,
				),
			)
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
	persons := []int{1, 2, 3}

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

			var ending string

			if number == "singular" {
				ending = "us"
			} else {
				ending = "ī"
			}

			form := ppp + ending + " " + sumForms[number][key]

			forms = append(forms,
				verbForm(
					form,
					person,
					number,
					tense,
					mood,
					"passive",
				),
			)
		}
	}

	return forms
}