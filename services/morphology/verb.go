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
		case 2:
			return generateSecondConjugation(lemma)
		case 3:
    		return generateThirdConjugation(lemma)
		case 31:
			return generateThirdIOConjugation(lemma)
		case 4:
			return generateFourthConjugation(lemma)
	}

	return []models.Form{}
}

func generateConjugation(
	lemma models.Lemma,
	infinitiveEnding string,
	passiveInfinitive string,
	gerundSuffix string,
	participleEnding string,
	patterns []FinitePattern,
	ppPatterns []PerfectPassivePattern,
	imperatives []ImperativePattern,
) []models.Form {

	var forms []models.Form

	// STEMS
	presentStem := removeVerbEnding(*lemma.Infinitive, infinitiveEnding)
	perfectStem := removeVerbEnding(*lemma.Perfect, "ī")
	ppp := pppStem(*lemma.Supine)

	gerundStem := presentStem + gerundSuffix
	gerundiveStem := buildGerundiveStem(lemma)

	pluperfectSubjStem := perfectStem + "isse"
	impSubPasStem := removeVerbEnding(*lemma.Infinitive, "re")

	// INFINITIVES
	forms = append(
		forms,
		buildInfinitives(
			presentStem,
			perfectStem,
			ppp,
			infinitiveEnding,
			passiveInfinitive,
		)...,
	)

	// GERUND
	forms = append(forms, buildGerundForms(gerundStem)...)

	// GERUNDIVE
	forms = append(forms, buildGerundiveForms(lemma, gerundiveStem)...)

	// PARTICIPLES
	forms = append(
		forms,
		generatePresentActiveParticiple(
			presentStem,
			participleEnding,
		)...,
	)

	forms = append(
		forms,
		generatePerfectPassiveParticiple(lemma, ppp)...,
	)

	forms = append(
		forms,
		generateFutureActiveParticiple(lemma, ppp)...,
	)

	// FINITE FORMS
	for _, pattern := range patterns {

		var stem string

		switch pattern.Stem {

		case "present":
			stem = presentStem

		case "perfect":
			stem = perfectStem

		case "pluperfect_subj":
			stem = pluperfectSubjStem

		case "imperfect_passive_subj":
			stem = impSubPasStem

		case "imperfect_subj":
			stem = *lemma.Infinitive

		default:
			panic("unknown stem type: " + pattern.Stem)
		}

		forms = append(
			forms,
			buildFiniteVerbForms(
				stem,
				pattern.Endings,
				pattern.Tense,
				pattern.Mood,
				pattern.Voice,
			)...,
		)
	}

	// PERFECT PASSIVE FORMS
	for _, pattern := range ppPatterns {

		forms = append(
			forms,
			buildPerfectPassiveForms(
				ppp,
				pattern.Endings,
				pattern.Tense,
				pattern.Mood,
			)...,
		)
	}

	// IMPERATIVES
	for _, pattern := range imperatives {

		forms = append(
			forms,
			imperativeForm(
				presentStem+pattern.Form,
				pattern.Person,
				pattern.Number,
				pattern.Tense,
				pattern.Voice,
			),
		)
	}

	return forms
}

func generateFirstConjugation(lemma models.Lemma) []models.Form {
	return generateConjugation(
		lemma,
		"āre",
		"ārī",
		"and",
		"āns",
		FirstConjugationPatterns,
		FirstConjugationPerfectPassivePatterns,
		FirstConjugationImperatives,
	)
}

func generateSecondConjugation(lemma models.Lemma) []models.Form {
	return generateConjugation(
		lemma,
		"ēre",
		"ērī",
		"end",
		"ēns",
		SecondConjugationPatterns,
		SecondConjugationPerfectPassivePatterns,
		SecondConjugationImperatives,
	)
}

func generateThirdConjugation(lemma models.Lemma) []models.Form {
	return generateConjugation(
		lemma,
		"ere",
		"ī",
		"end",
		"ēns",
		ThirdConjugationPatterns,
		ThirdConjugationPerfectPassivePatterns,
		ThirdConjugationImperatives,
	)
}

func generateThirdIOConjugation(lemma models.Lemma) []models.Form {
	return generateConjugation(
		lemma,
		"ere",
		"ī",
		"iend",
		"iēns",
		ThirdIOConjugationPatterns,
		ThirdIOConjugationPerfectPassivePatterns,
		ThirdIOConjugationImperatives,
	)
}

func generateFourthConjugation(lemma models.Lemma) []models.Form {
	return generateConjugation(
		lemma,
		"īre",   // infinitive ending
		"īrī",   // passive infinitive
		"iend",  // gerund suffix
		"iēns",  // present active participle
		FourthConjugationPatterns,
		FourthConjugationPerfectPassivePatterns,
		FourthConjugationImperatives,
	)
}