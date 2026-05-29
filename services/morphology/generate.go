package morphology

import "iuno-api/models"

func Generate(lemma models.Lemma,) []models.Form {

	// REGULAR SYSTEMS
	switch lemma.PartOfSpeech {

		case "noun":
			return GenerateNoun(lemma)

		case "verb":
			return GenerateVerb(lemma)

		case "adjective":
			return GenerateAdjective(lemma)

	}

	return []models.Form{}
}