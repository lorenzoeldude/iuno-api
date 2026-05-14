package morphology

// import "iuno-api/models"

//
// VERB ENDING SYSTEM
//
// This file contains ONLY Latin endings.
//
// No conjugation logic belongs here.
//
// The morphology engine combines:
//
// stem + ending -> final form
//

// =====================================================
// PERSON KEYS
// =====================================================
//
// 1s = 1st singular
// 2s = 2nd singular
// 3s = 3rd singular
// 1p = 1st plural
// 2p = 2nd plural
// 3p = 3rd plural
//

// =====================================================
// PRESENT ACTIVE INDICATIVE
// =====================================================

var PresentActiveEndings = map[int]map[string]string{

	// 1st conjugation
	1: {
		"1s": "o",
		"2s": "as",
		"3s": "at",
		"1p": "amus",
		"2p": "atis",
		"3p": "ant",
	},

	// 2nd conjugation
	2: {
		"1s": "eo",
		"2s": "es",
		"3s": "et",
		"1p": "emus",
		"2p": "etis",
		"3p": "ent",
	},

	// 3rd conjugation
	3: {
		"1s": "o",
		"2s": "is",
		"3s": "it",
		"1p": "imus",
		"2p": "itis",
		"3p": "unt",
	},

	// 4th conjugation
	4: {
		"1s": "io",
		"2s": "is",
		"3s": "it",
		"1p": "imus",
		"2p": "itis",
		"3p": "iunt",
	},
}

// =====================================================
// IMPERFECT ACTIVE INDICATIVE
// =====================================================

var ImperfectActiveEndings = map[int]map[string]string{

	1: {
		"1s": "abam",
		"2s": "abas",
		"3s": "abat",
		"1p": "abamus",
		"2p": "abatis",
		"3p": "abant",
	},

	2: {
		"1s": "ebam",
		"2s": "ebas",
		"3s": "ebat",
		"1p": "ebamus",
		"2p": "ebatis",
		"3p": "ebant",
	},

	3: {
		"1s": "ebam",
		"2s": "ebas",
		"3s": "ebat",
		"1p": "ebamus",
		"2p": "ebatis",
		"3p": "ebant",
	},

	4: {
		"1s": "iebam",
		"2s": "iebas",
		"3s": "iebat",
		"1p": "iebamus",
		"2p": "iebatis",
		"3p": "iebant",
	},
}

// =====================================================
// FUTURE ACTIVE INDICATIVE
// =====================================================

var FutureActiveEndings = map[int]map[string]string{

	1: {
		"1s": "abo",
		"2s": "abis",
		"3s": "abit",
		"1p": "abimus",
		"2p": "abitis",
		"3p": "abunt",
	},

	2: {
		"1s": "ebo",
		"2s": "ebis",
		"3s": "ebit",
		"1p": "ebimus",
		"2p": "ebitis",
		"3p": "ebunt",
	},

	3: {
		"1s": "am",
		"2s": "es",
		"3s": "et",
		"1p": "emus",
		"2p": "etis",
		"3p": "ent",
	},

	4: {
		"1s": "iam",
		"2s": "ies",
		"3s": "iet",
		"1p": "iemus",
		"2p": "ietis",
		"3p": "ient",
	},
}

// =====================================================
// PERFECT ACTIVE INDICATIVE
// =====================================================

var PerfectActiveEndings = map[string]string{
	"1s": "i",
	"2s": "isti",
	"3s": "it",
	"1p": "imus",
	"2p": "istis",
	"3p": "erunt",
}

// =====================================================
// PLUPERFECT ACTIVE INDICATIVE
// =====================================================

var PluperfectActiveEndings = map[string]string{
	"1s": "eram",
	"2s": "eras",
	"3s": "erat",
	"1p": "eramus",
	"2p": "eratis",
	"3p": "erant",
}

// =====================================================
// FUTURE PERFECT ACTIVE INDICATIVE
// =====================================================

var FuturePerfectActiveEndings = map[string]string{
	"1s": "ero",
	"2s": "eris",
	"3s": "erit",
	"1p": "erimus",
	"2p": "eritis",
	"3p": "erint",
}

// =====================================================
// PERSONAL ENDING HELPER
// =====================================================

func personKey(
	person int,
	number string,
) string {

	if number == "singular" {

		switch person {
		case 1:
			return "1s"
		case 2:
			return "2s"
		case 3:
			return "3s"
		}
	}

	if number == "plural" {

		switch person {
		case 1:
			return "1p"
		case 2:
			return "2p"
		case 3:
			return "3p"
		}
	}

	return ""
}