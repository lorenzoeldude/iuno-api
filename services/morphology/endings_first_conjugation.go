package morphology

import (
	// "iuno-api/models"
)

type VerbEndings map[string]map[string]string

type FinitePattern struct {
    Stem    string
    Endings map[string]map[string]string
    Tense   string
    Mood    string
    Voice   string
}

type PerfectPassivePattern struct {
    Endings VerbEndings
    Tense   string
    Mood    string
}

type ImperativePattern struct {
    Form   string
    Person int
    Number string
    Tense  string
    Voice  string
}


var firstPresentActiveIndicative = VerbEndings{
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
}

var firstImperfectActiveIndicative = VerbEndings{
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
}

var firstFutureActiveIndicative = VerbEndings{
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
}

var firstPerfectActiveIndicative = VerbEndings{
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
}

var firstPluperfectActiveIndicative = VerbEndings{
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
}

var firstFuturePerfectActiveIndicative = VerbEndings{
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
}

var firstPresentActiveSubjunctive = VerbEndings{
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
}

var firstImperfectActiveSubjunctive = VerbEndings{
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
}

var firstPerfectActiveSubjunctive = VerbEndings{
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
}

var firstPluperfectActiveSubjunctive = VerbEndings{
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
}

var firstPresentPassiveIndicative = VerbEndings{
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
}

var firstImperfectPassiveIndicative = VerbEndings{
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
}

var firstFuturePassiveIndicative = VerbEndings{
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
}

var firstPerfectPassiveIndicative = VerbEndings{
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
}

var firstPluperfectPassiveIndicative = VerbEndings{
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
}

var firstFuturePerfectPassiveIndicative = VerbEndings{
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
}

var firstPresentPassiveSubjunctive = VerbEndings{
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
}

var firstImperfectPassiveSubjunctive = VerbEndings{
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
}

var firstPerfectPassiveSubjunctive = VerbEndings{
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
}

var firstPluperfectPassiveSubjunctive = VerbEndings{
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
}

var FirstConjugationImperatives = []ImperativePattern{
    // Present Active
    {"ā", 2, "singular", "present", "active"},
    {"āte", 2, "plural", "present", "active"},

    // Future Active
    {"ātō", 2, "singular", "future", "active"},
    {"ātōte", 2, "plural", "future", "active"},
    {"ātō", 3, "singular", "future", "active"},
    {"antō", 3, "plural", "future", "active"},

    // Present Passive
    {"āre", 2, "singular", "present", "passive"},
    {"āminī", 2, "plural", "present", "passive"},

    // Future Passive
    {"ātor", 2, "singular", "future", "passive"},
    {"āminī", 2, "plural", "future", "passive"},
    {"ātor", 3, "singular", "future", "passive"},
    {"antor", 3, "plural", "future", "passive"},
}

var FirstConjugationPatterns = []FinitePattern{

    // ACTIVE INDICATIVE

    {
        Stem:    "present",
        Endings: firstPresentActiveIndicative,
        Tense:   "present",
        Mood:    "indicative",
        Voice:   "active",
    },
    {
        Stem:    "present",
        Endings: firstImperfectActiveIndicative,
        Tense:   "imperfect",
        Mood:    "indicative",
        Voice:   "active",
    },
    {
        Stem:    "present",
        Endings: firstFutureActiveIndicative,
        Tense:   "future",
        Mood:    "indicative",
        Voice:   "active",
    },
    {
        Stem:    "perfect",
        Endings: firstPerfectActiveIndicative,
        Tense:   "perfect",
        Mood:    "indicative",
        Voice:   "active",
    },
    {
        Stem:    "perfect",
        Endings: firstPluperfectActiveIndicative,
        Tense:   "pluperfect",
        Mood:    "indicative",
        Voice:   "active",
    },
    {
        Stem:    "perfect",
        Endings: firstFuturePerfectActiveIndicative,
        Tense:   "future perfect",
        Mood:    "indicative",
        Voice:   "active",
    },

    // ACTIVE SUBJUNCTIVE

    {
        Stem:    "present",
        Endings: firstPresentActiveSubjunctive,
        Tense:   "present",
        Mood:    "subjunctive",
        Voice:   "active",
    },
    {
        Stem:    "imperfect_subj",
        Endings: firstImperfectActiveSubjunctive,
        Tense:   "imperfect",
        Mood:    "subjunctive",
        Voice:   "active",
    },
    {
        Stem:    "perfect",
        Endings: firstPerfectActiveSubjunctive,
        Tense:   "perfect",
        Mood:    "subjunctive",
        Voice:   "active",
    },
    {
        Stem:    "pluperfect_subj",
        Endings: firstPluperfectActiveSubjunctive,
        Tense:   "pluperfect",
        Mood:    "subjunctive",
        Voice:   "active",
    },

    // PASSIVE INDICATIVE

    {
        Stem:    "present",
        Endings: firstPresentPassiveIndicative,
        Tense:   "present",
        Mood:    "indicative",
        Voice:   "passive",
    },
    {
        Stem:    "present",
        Endings: firstImperfectPassiveIndicative,
        Tense:   "imperfect",
        Mood:    "indicative",
        Voice:   "passive",
    },
    {
        Stem:    "present",
        Endings: firstFuturePassiveIndicative,
        Tense:   "future",
        Mood:    "indicative",
        Voice:   "passive",
    },

    // PASSIVE SUBJUNCTIVE

    {
        Stem:    "present",
        Endings: firstPresentPassiveSubjunctive,
        Tense:   "present",
        Mood:    "subjunctive",
        Voice:   "passive",
    },
    {
        Stem:    "imperfect_passive_subj",
        Endings: firstImperfectPassiveSubjunctive,
        Tense:   "imperfect",
        Mood:    "subjunctive",
        Voice:   "passive",
    },
}

var FirstConjugationPerfectPassivePatterns = []PerfectPassivePattern{
    {
        Endings: firstPerfectPassiveIndicative,
        Tense:   "perfect",
        Mood:    "indicative",
    },
    {
        Endings: firstPluperfectPassiveIndicative,
        Tense:   "pluperfect",
        Mood:    "indicative",
    },
    {
        Endings: firstFuturePerfectPassiveIndicative,
        Tense:   "future perfect",
        Mood:    "indicative",
    },
    {
        Endings: firstPerfectPassiveSubjunctive,
        Tense:   "perfect",
        Mood:    "subjunctive",
    },
    {
        Endings: firstPluperfectPassiveSubjunctive,
        Tense:   "pluperfect",
        Mood:    "subjunctive",
    },
}