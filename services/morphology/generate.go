package morphology

import "iuno-api/models"

func Generate(
	word models.Word,
) []models.Form {

	// =====================================================
	// IRREGULAR OVERRIDE
	// =====================================================

	if word.Irregular {
		return getIrregularForms(word)
	}

	// =====================================================
	// REGULAR SYSTEMS
	// =====================================================

	switch word.Type {

	case "noun":
		return GenerateNoun(word)

	case "verb":
		return GenerateVerb(word)

	case "adjective":
		return GenerateAdjective(word)

	case "pronoun":
		return GeneratePronoun(word)
	}

	return []models.Form{}
}