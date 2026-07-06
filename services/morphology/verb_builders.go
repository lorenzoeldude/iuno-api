package morphology

import (
	"iuno-api/models"
	"strings"
)

// INFINITIVES (REUSABLE)
func buildInfinitives(
	presentStem string,
	perfectStem string,
	ppp string,
	presentActiveEnding string,
	presentPassiveEnding string,
) []models.Form {

	var forms []models.Form

	// Always available
	forms = append(forms,
		infinitiveForm(presentStem+presentActiveEnding, "present", "active"),
		infinitiveForm(perfectStem+"isse", "perfect", "active"),
		infinitiveForm(presentStem+presentPassiveEnding, "present", "passive"),
	)

	// Only if a supine exists
	if ppp != "" {
		forms = append(forms,
			infinitiveForm(ppp+"ūrus esse", "future", "active"),
			infinitiveForm(ppp+"um esse", "perfect", "passive"),
			infinitiveForm(ppp+"um īrī", "future", "passive"),
		)
	}

	return forms
}

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

// GERUND (REUSABLE)
func buildGerundForms(gerundStem string) []models.Form {

	return []models.Form{
		{
			Form:         gerundStem + "ī",
			PartOfSpeech: "verb",
			GrammaticalCase: StringPtr("genitive"),
			Number:       "singular",
			Mood:         StringPtr("gerund"),
		},
		{
			Form:         gerundStem + "ō",
			PartOfSpeech: "verb",
			GrammaticalCase: StringPtr("dative"),
			Number:       "singular",
			Mood:         StringPtr("gerund"),
		},
		{
			Form:         gerundStem + "um",
			PartOfSpeech: "verb",
			GrammaticalCase: StringPtr("accusative"),
			Number:       "singular",
			Mood:         StringPtr("gerund"),
		},
		{
			Form:         gerundStem + "ō",
			PartOfSpeech: "verb",
			GrammaticalCase: StringPtr("ablative"),
			Number:       "singular",
			Mood:         StringPtr("gerund"),
		},
	}
}

// 	GERUNDIVE (REUSABLE)
func buildGerundiveForms(lemma models.Lemma, gerundiveStem string) []models.Form {
	var gerundiveForms []models.Form

	gerundiveForms = append(
		gerundiveForms,
		buildMasculineAdjectiveForms(lemma, gerundiveStem)...,
	)

	gerundiveForms = append(
		gerundiveForms,
		buildFeminineAdjectiveForms(lemma, gerundiveStem)...,
	)

	gerundiveForms = append(
		gerundiveForms,
		buildNeuterAdjectiveForms(lemma, gerundiveStem)...,
	)

	mood := "gerundive"

	for i := range gerundiveForms {
		gerundiveForms[i].PartOfSpeech = "verb"
		gerundiveForms[i].Mood = &mood
	}

	return gerundiveForms
}

// PARTICIPLES
func markAsParticiple(
    forms []models.Form,
    tense string,
    voice string,
) {
    for i := range forms {
        forms[i].PartOfSpeech = "verb"
        forms[i].Mood = StringPtr("participle")
        forms[i].Tense = StringPtr(tense)
        forms[i].Voice = StringPtr(voice)
    }
}

// PAP (REUSABLE)
func generatePresentActiveParticiple(
	presentStem string,
    papEnding string,
) []models.Form {

    var forms []models.Form

	nominativeForm := presentStem + papEnding

	papStem := strings.TrimSuffix(NormalizeLatin(nominativeForm), "s")


    endings := map[string]map[string]string{
        "singular": {
            "genitive":   "tis",
            "dative":     "tī",
            "accusative": "tem", // m/f only
            "ablative":   "te",  // participles usually use -e
        },

        "plural": {
            "nominative": "tēs", // m/f only
            "genitive":   "tium",
            "dative":     "tibus",
            "accusative": "tēs", // m/f only
            "ablative":   "tibus",
            "vocative":   "tēs", // m/f only
        },
    }

    papLemma := models.Lemma{
		Declension: IntPtr(31),
		Lemma:      nominativeForm,
		Neuter:     StringPtr(papStem + "s"),
	}

	for _, gender := range []string{
		"masculine",
		"feminine",
		"neuter",
	} {
		forms = append(
			forms,
			buildThirdDeclensionForms(
				papLemma,
				papStem,
				gender,
				endings,
			)...,
		)
	}

    markAsParticiple(
        forms,
        "present",
        "active",
    )

    return forms
}

// PPP (REUSABLE)
func generatePerfectPassiveParticiple(lemma models.Lemma, ppp string) []models.Form {
	
	var pppForms []models.Form
	
	pppForms = append(
		pppForms,
		buildMasculineAdjectiveForms(lemma, ppp)...,
	)

	pppForms = append(
		pppForms,
		buildFeminineAdjectiveForms(lemma, ppp)...,
	)

	pppForms = append(
		pppForms,
		buildNeuterAdjectiveForms(lemma, ppp)...,
	)

	markAsParticiple(pppForms, "perfect", "passive",)

	return pppForms
}

// FAP (REUSABLE)
func generateFutureActiveParticiple(lemma models.Lemma, ppp string) []models.Form {
	var fapForms []models.Form

	futureActiveStem := ppp + "ūr"

	fapForms = append(
		fapForms,
		buildMasculineAdjectiveForms(lemma, futureActiveStem)...,
	)

	fapForms = append(
		fapForms,
		buildFeminineAdjectiveForms(lemma, futureActiveStem)...,
	)

	fapForms = append(
		fapForms,
		buildNeuterAdjectiveForms(lemma, futureActiveStem)...,
	)

	markAsParticiple(
		fapForms,
		"future",
		"active",
	)

	return fapForms
}

func imperativeForm(
    form string,
    person int,
    number string,
    tense string,
    voice string,
) models.Form {
    return models.Form{
        Form:         form,
        PartOfSpeech: "verb",
        Person:       IntPtr(person),
        Number:       number,
        Tense:        StringPtr(tense),
        Mood:         StringPtr("imperative"),
        Voice:        StringPtr(voice),
    }
}

// BUILD FINITE VERB FORMS
func buildFiniteVerbForms(
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

// BUILD PERFECT PASSIVE FORMS
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