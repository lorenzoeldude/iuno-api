package morphology

import "iuno-api/models"

func Generate(lemma models.Lemma,) []models.Form {

	// IRREGULAR OVERRIDE
	if lemma.Irregular {
		return getIrregularForms(lemma)
	}

	// REGULAR SYSTEMS
	switch lemma.Type {

		case "noun":
			return GenerateNoun(lemma)

		// case "verb":
		// 	return GenerateVerb(lemma)

		case "adjective":
			return GenerateAdjective(lemma)

		case "pronoun":
			return GeneratePronoun(lemma)
	}

	return []models.Form{}
}