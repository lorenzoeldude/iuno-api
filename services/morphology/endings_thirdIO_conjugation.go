package morphology

var thirdIOPresentActiveIndicative = VerbEndings{
	"singular": {
		"first":  "iō",
		"second": "is",
		"third":  "it",
	},
	"plural": {
		"first":  "imus",
		"second": "itis",
		"third":  "iunt",
	},
}

var thirdIOImperfectActiveIndicative = VerbEndings{
	"singular": {
		"first":  "iēbam",
		"second": "iēbās",
		"third":  "iēbat",
	},
	"plural": {
		"first":  "iēbāmus",
		"second": "iēbātis",
		"third":  "iēbant",
	},
}

var thirdIOFutureActiveIndicative = VerbEndings{
	"singular": {
		"first":  "iam",
		"second": "iēs",
		"third":  "iet",
	},
	"plural": {
		"first":  "iēmus",
		"second": "iētis",
		"third":  "ient",
	},
}

var thirdIOPresentActiveSubjunctive = VerbEndings{
	"singular": {
		"first":  "iam",
		"second": "iās",
		"third":  "iat",
	},
	"plural": {
		"first":  "iāmus",
		"second": "iātis",
		"third":  "iant",
	},
}

var thirdIOPresentPassiveIndicative = VerbEndings{
	"singular": {
		"first":  "ior",
		"second": "eris",
		"third":  "itur",
	},
	"plural": {
		"first":  "imur",
		"second": "iminī",
		"third":  "iuntur",
	},
}

var thirdIOImperfectPassiveIndicative = VerbEndings{
	"singular": {
		"first":  "iēbar",
		"second": "iēbāris",
		"third":  "iēbātur",
	},
	"plural": {
		"first":  "iēbāmur",
		"second": "iēbāminī",
		"third":  "iēbantur",
	},
}

var thirdIOFuturePassiveIndicative = VerbEndings{
	"singular": {
		"first":  "iar",
		"second": "iēris",
		"third":  "iētur",
	},
	"plural": {
		"first":  "iēmur",
		"second": "iēminī",
		"third":  "ientur",
	},
}

var thirdIOPresentPassiveSubjunctive = VerbEndings{
	"singular": {
		"first":  "iar",
		"second": "iāris",
		"third":  "iātur",
	},
	"plural": {
		"first":  "iāmur",
		"second": "iāminī",
		"third":  "iantur",
	},
}

var thirdIOPerfectActiveIndicative       = firstPerfectActiveIndicative
var thirdIOPluperfectActiveIndicative    = firstPluperfectActiveIndicative
var thirdIOFuturePerfectActiveIndicative = firstFuturePerfectActiveIndicative

var thirdIOImperfectActiveSubjunctive    = firstImperfectActiveSubjunctive
var thirdIOPerfectActiveSubjunctive      = firstPerfectActiveSubjunctive
var thirdIOPluperfectActiveSubjunctive   = firstPluperfectActiveSubjunctive

var thirdIOPerfectPassiveIndicative       = firstPerfectPassiveIndicative
var thirdIOPluperfectPassiveIndicative    = firstPluperfectPassiveIndicative
var thirdIOFuturePerfectPassiveIndicative = firstFuturePerfectPassiveIndicative

var thirdIOImperfectPassiveSubjunctive    = firstImperfectPassiveSubjunctive
var thirdIOPerfectPassiveSubjunctive      = firstPerfectPassiveSubjunctive
var thirdIOPluperfectPassiveSubjunctive   = firstPluperfectPassiveSubjunctive

var ThirdIOConjugationImperatives = []ImperativePattern{

	// Present Active
	{"e", 2, "singular", "present", "active"},
	{"ite", 2, "plural", "present", "active"},

	// Future Active
	{"itō", 2, "singular", "future", "active"},
	{"itōte", 2, "plural", "future", "active"},
	{"itō", 3, "singular", "future", "active"},
	{"iuntō", 3, "plural", "future", "active"},

	// Present Passive
	{"ere", 2, "singular", "present", "passive"},
	{"iminī", 2, "plural", "present", "passive"},

	// Future Passive
	{"itor", 2, "singular", "future", "passive"},
	{"iminī", 2, "plural", "future", "passive"},
	{"itor", 3, "singular", "future", "passive"},
	{"iuntor", 3, "plural", "future", "passive"},
}

var ThirdIOConjugationPatterns = []FinitePattern{

	// ACTIVE INDICATIVE

	{
		Stem:    "present",
		Endings: thirdIOPresentActiveIndicative,
		Tense:   "present",
		Mood:    "indicative",
		Voice:   "active",
	},
	{
		Stem:    "present",
		Endings: thirdIOImperfectActiveIndicative,
		Tense:   "imperfect",
		Mood:    "indicative",
		Voice:   "active",
	},
	{
		Stem:    "present",
		Endings: thirdIOFutureActiveIndicative,
		Tense:   "future",
		Mood:    "indicative",
		Voice:   "active",
	},
	{
		Stem:    "perfect",
		Endings: thirdPerfectActiveIndicative,
		Tense:   "perfect",
		Mood:    "indicative",
		Voice:   "active",
	},
	{
		Stem:    "perfect",
		Endings: thirdPluperfectActiveIndicative,
		Tense:   "pluperfect",
		Mood:    "indicative",
		Voice:   "active",
	},
	{
		Stem:    "perfect",
		Endings: thirdFuturePerfectActiveIndicative,
		Tense:   "future perfect",
		Mood:    "indicative",
		Voice:   "active",
	},

	// ACTIVE SUBJUNCTIVE

	{
		Stem:    "present",
		Endings: thirdIOPresentActiveSubjunctive,
		Tense:   "present",
		Mood:    "subjunctive",
		Voice:   "active",
	},
	{
		Stem:    "imperfect_subj",
		Endings: thirdImperfectActiveSubjunctive,
		Tense:   "imperfect",
		Mood:    "subjunctive",
		Voice:   "active",
	},
	{
		Stem:    "perfect",
		Endings: thirdPerfectActiveSubjunctive,
		Tense:   "perfect",
		Mood:    "subjunctive",
		Voice:   "active",
	},
	{
		Stem:    "pluperfect_subj",
		Endings: thirdPluperfectActiveSubjunctive,
		Tense:   "pluperfect",
		Mood:    "subjunctive",
		Voice:   "active",
	},

	// PASSIVE INDICATIVE

	{
		Stem:    "present",
		Endings: thirdIOPresentPassiveIndicative,
		Tense:   "present",
		Mood:    "indicative",
		Voice:   "passive",
	},
	{
		Stem:    "present",
		Endings: thirdIOImperfectPassiveIndicative,
		Tense:   "imperfect",
		Mood:    "indicative",
		Voice:   "passive",
	},
	{
		Stem:    "present",
		Endings: thirdIOFuturePassiveIndicative,
		Tense:   "future",
		Mood:    "indicative",
		Voice:   "passive",
	},

	// PASSIVE SUBJUNCTIVE

	{
		Stem:    "present",
		Endings: thirdIOPresentPassiveSubjunctive,
		Tense:   "present",
		Mood:    "subjunctive",
		Voice:   "passive",
	},
	{
		Stem:    "imperfect_passive_subj",
		Endings: thirdImperfectPassiveSubjunctive,
		Tense:   "imperfect",
		Mood:    "subjunctive",
		Voice:   "passive",
	},
}

var ThirdIOConjugationPerfectPassivePatterns = ThirdConjugationPerfectPassivePatterns