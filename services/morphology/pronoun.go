package morphology

import "iuno-api/models"

//
// PRONOUN MORPHOLOGY ENGINE
//

func GeneratePronoun(word models.Word) []models.Form {

	switch word.Lemma {

	case "ego":
		return generateEgo()

	case "tu":
		return generateTu()

	case "nos":
		return generateNos()

	case "vos":
		return generateVos()

	case "is":
		return generateIsEaId()

	case "hic":
		return generateHicHaecHoc()

	case "ille":
		return generateIlleIllaIllud()

	case "qui":
		return generateQuiQuaeQuod()
	}

	return []models.Form{}
}

// =====================================================
// EGO
// =====================================================

func generateEgo() []models.Form {

	return []models.Form{

		{
			Form: "ego",
			Part: "pronoun",
			Case: "nominative",
			Number: "singular",
			Person: 1,
		},

		{
			Form: "mei",
			Part: "pronoun",
			Case: "genitive",
			Number: "singular",
			Person: 1,
		},

		{
			Form: "mihi",
			Part: "pronoun",
			Case: "dative",
			Number: "singular",
			Person: 1,
		},

		{
			Form: "me",
			Part: "pronoun",
			Case: "accusative",
			Number: "singular",
			Person: 1,
		},

		{
			Form: "me",
			Part: "pronoun",
			Case: "ablative",
			Number: "singular",
			Person: 1,
		},
	}
}

// =====================================================
// TU
// =====================================================

func generateTu() []models.Form {

	return []models.Form{

		{
			Form: "tu",
			Part: "pronoun",
			Case: "nominative",
			Number: "singular",
			Person: 2,
		},

		{
			Form: "tui",
			Part: "pronoun",
			Case: "genitive",
			Number: "singular",
			Person: 2,
		},

		{
			Form: "tibi",
			Part: "pronoun",
			Case: "dative",
			Number: "singular",
			Person: 2,
		},

		{
			Form: "te",
			Part: "pronoun",
			Case: "accusative",
			Number: "singular",
			Person: 2,
		},

		{
			Form: "te",
			Part: "pronoun",
			Case: "ablative",
			Number: "singular",
			Person: 2,
		},
	}
}

// =====================================================
// NOS
// =====================================================

func generateNos() []models.Form {

	return []models.Form{

		{
			Form: "nos",
			Part: "pronoun",
			Case: "nominative",
			Number: "plural",
			Person: 1,
		},

		{
			Form: "nostri",
			Part: "pronoun",
			Case: "genitive",
			Number: "plural",
			Person: 1,
		},

		{
			Form: "nobis",
			Part: "pronoun",
			Case: "dative",
			Number: "plural",
			Person: 1,
		},

		{
			Form: "nos",
			Part: "pronoun",
			Case: "accusative",
			Number: "plural",
			Person: 1,
		},

		{
			Form: "nobis",
			Part: "pronoun",
			Case: "ablative",
			Number: "plural",
			Person: 1,
		},
	}
}

// =====================================================
// VOS
// =====================================================

func generateVos() []models.Form {

	return []models.Form{

		{
			Form: "vos",
			Part: "pronoun",
			Case: "nominative",
			Number: "plural",
			Person: 2,
		},

		{
			Form: "vestri",
			Part: "pronoun",
			Case: "genitive",
			Number: "plural",
			Person: 2,
		},

		{
			Form: "vobis",
			Part: "pronoun",
			Case: "dative",
			Number: "plural",
			Person: 2,
		},

		{
			Form: "vos",
			Part: "pronoun",
			Case: "accusative",
			Number: "plural",
			Person: 2,
		},

		{
			Form: "vobis",
			Part: "pronoun",
			Case: "ablative",
			Number: "plural",
			Person: 2,
		},
	}
}

// =====================================================
// IS EA ID
// =====================================================

func generateIsEaId() []models.Form {

	return []models.Form{

		{
			Form: "is",
			Part: "pronoun",
			Case: "nominative",
			Number: "singular",
			Gender: "masculine",
		},

		{
			Form: "eius",
			Part: "pronoun",
			Case: "genitive",
			Number: "singular",
			Gender: "masculine",
		},

		{
			Form: "ei",
			Part: "pronoun",
			Case: "dative",
			Number: "singular",
			Gender: "masculine",
		},

		{
			Form: "eum",
			Part: "pronoun",
			Case: "accusative",
			Number: "singular",
			Gender: "masculine",
		},

		{
			Form: "eo",
			Part: "pronoun",
			Case: "ablative",
			Number: "singular",
			Gender: "masculine",
		},

		{
			Form: "ea",
			Part: "pronoun",
			Case: "nominative",
			Number: "singular",
			Gender: "feminine",
		},

		{
			Form: "eam",
			Part: "pronoun",
			Case: "accusative",
			Number: "singular",
			Gender: "feminine",
		},

		{
			Form: "id",
			Part: "pronoun",
			Case: "nominative",
			Number: "singular",
			Gender: "neuter",
		},

		{
			Form: "id",
			Part: "pronoun",
			Case: "accusative",
			Number: "singular",
			Gender: "neuter",
		},
	}
}

// =====================================================
// HIC HAEC HOC
// =====================================================

func generateHicHaecHoc() []models.Form {

	return []models.Form{

		{
			Form: "hic",
			Part: "pronoun",
			Case: "nominative",
			Number: "singular",
			Gender: "masculine",
		},

		{
			Form: "huius",
			Part: "pronoun",
			Case: "genitive",
			Number: "singular",
			Gender: "masculine",
		},

		{
			Form: "huic",
			Part: "pronoun",
			Case: "dative",
			Number: "singular",
			Gender: "masculine",
		},

		{
			Form: "hunc",
			Part: "pronoun",
			Case: "accusative",
			Number: "singular",
			Gender: "masculine",
		},

		{
			Form: "hoc",
			Part: "pronoun",
			Case: "ablative",
			Number: "singular",
			Gender: "masculine",
		},
	}
}

// =====================================================
// ILLE ILLA ILLUD
// =====================================================

func generateIlleIllaIllud() []models.Form {

	return []models.Form{

		{
			Form: "ille",
			Part: "pronoun",
			Case: "nominative",
			Number: "singular",
			Gender: "masculine",
		},

		{
			Form: "illum",
			Part: "pronoun",
			Case: "accusative",
			Number: "singular",
			Gender: "masculine",
		},

		{
			Form: "illa",
			Part: "pronoun",
			Case: "nominative",
			Number: "singular",
			Gender: "feminine",
		},

		{
			Form: "illud",
			Part: "pronoun",
			Case: "nominative",
			Number: "singular",
			Gender: "neuter",
		},
	}
}

// =====================================================
// QUI QUAE QUOD
// =====================================================

func generateQuiQuaeQuod() []models.Form {

	return []models.Form{

		{
			Form: "qui",
			Part: "pronoun",
			Case: "nominative",
			Number: "singular",
			Gender: "masculine",
		},

		{
			Form: "quem",
			Part: "pronoun",
			Case: "accusative",
			Number: "singular",
			Gender: "masculine",
		},

		{
			Form: "quae",
			Part: "pronoun",
			Case: "nominative",
			Number: "singular",
			Gender: "feminine",
		},

		{
			Form: "quod",
			Part: "pronoun",
			Case: "nominative",
			Number: "singular",
			Gender: "neuter",
		},
	}
}