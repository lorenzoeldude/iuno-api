package morphology

import "iuno-api/models"

//
// PREPOSITIONS
//
// Latin prepositions are mostly indeclinable.
// They do not generate paradigms like nouns or verbs.
//
// But they DO carry grammatical metadata,
// especially which case they govern.
//
// Examples:
//
// ad   + accusative
// cum  + ablative
// in   + ablative/accusative
//

func GeneratePreposition(
	word models.Word,
) []models.Form {

	var governedCase string

	switch word.Lemma {

	// =====================================================
	// ACCUSATIVE
	// =====================================================

	case "ad":
		governedCase = "accusative"

	case "per":
		governedCase = "accusative"

	case "propter":
		governedCase = "accusative"

	case "contra":
		governedCase = "accusative"

	case "post":
		governedCase = "accusative"

	case "ante":
		governedCase = "accusative"

	case "inter":
		governedCase = "accusative"

	case "trans":
		governedCase = "accusative"

	case "circum":
		governedCase = "accusative"

	// =====================================================
	// ABLATIVE
	// =====================================================

	case "cum":
		governedCase = "ablative"

	case "de":
		governedCase = "ablative"

	case "ex":
		governedCase = "ablative"

	case "e":
		governedCase = "ablative"

	case "pro":
		governedCase = "ablative"

	case "sine":
		governedCase = "ablative"

	case "ab":
		governedCase = "ablative"

	case "a":
		governedCase = "ablative"

	// =====================================================
	// BOTH
	// =====================================================

	case "in":
		governedCase = "ablative/accusative"

	case "sub":
		governedCase = "ablative/accusative"

	case "super":
		governedCase = "ablative/accusative"

	// =====================================================
	// UNKNOWN
	// =====================================================

	default:
		governedCase = ""
	}

	return []models.Form{
		{
			Form:   word.Lemma,
			Part:   "preposition",
			Case:   governedCase,
			Number: "",
			Gender: "",

			Tense: "",
			Mood:  "",
			Voice: "",

			Person: 0,

			NonFinite: "",
		},
	}
}