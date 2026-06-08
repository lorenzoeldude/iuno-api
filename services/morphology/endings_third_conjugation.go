package morphology

var thirdPresentActiveIndicative = VerbEndings{
	"singular": {
		"first":  "ō",
		"second": "is",
		"third":  "it",
	},
	"plural": {
		"first":  "imus",
		"second": "itis",
		"third":  "unt",
	},
}

var thirdImperfectActiveIndicative = VerbEndings{
	"singular": {
		"first":  "ēbam",
		"second": "ēbās",
		"third":  "ēbat",
	},
	"plural": {
		"first":  "ēbāmus",
		"second": "ēbātis",
		"third":  "ēbant",
	},
}

var thirdFutureActiveIndicative = VerbEndings{
	"singular": {
		"first":  "am",
		"second": "ēs",
		"third":  "et",
	},
	"plural": {
		"first":  "ēmus",
		"second": "ētis",
		"third":  "ent",
	},
}

var thirdPresentActiveSubjunctive = VerbEndings{
	"singular": {
		"first":  "am",
		"second": "ās",
		"third":  "at",
	},
	"plural": {
		"first":  "āmus",
		"second": "ātis",
		"third":  "ant",
	},
}

var thirdPresentPassiveIndicative = VerbEndings{
	"singular": {
		"first":  "or",
		"second": "eris",
		"third":  "itur",
	},
	"plural": {
		"first":  "imur",
		"second": "iminī",
		"third":  "untur",
	},
}

var thirdImperfectPassiveIndicative = VerbEndings{
	"singular": {
		"first":  "ēbar",
		"second": "ēbāris",
		"third":  "ēbātur",
	},
	"plural": {
		"first":  "ēbāmur",
		"second": "ēbāminī",
		"third":  "ēbantur",
	},
}

var thirdFuturePassiveIndicative = VerbEndings{
	"singular": {
		"first":  "ar",
		"second": "ēris",
		"third":  "ētur",
	},
	"plural": {
		"first":  "ēmur",
		"second": "ēminī",
		"third":  "entur",
	},
}

var thirdPresentPassiveSubjunctive = VerbEndings{
	"singular": {
		"first":  "ar",
		"second": "āris",
		"third":  "ātur",
	},
	"plural": {
		"first":  "āmur",
		"second": "āminī",
		"third":  "antur",
	},
}

var thirdPerfectActiveIndicative       = firstPerfectActiveIndicative
var thirdPluperfectActiveIndicative    = firstPluperfectActiveIndicative
var thirdFuturePerfectActiveIndicative = firstFuturePerfectActiveIndicative

var thirdImperfectActiveSubjunctive    = firstImperfectActiveSubjunctive
var thirdPerfectActiveSubjunctive      = firstPerfectActiveSubjunctive
var thirdPluperfectActiveSubjunctive   = firstPluperfectActiveSubjunctive

var thirdPerfectPassiveIndicative       = firstPerfectPassiveIndicative
var thirdPluperfectPassiveIndicative    = firstPluperfectPassiveIndicative
var thirdFuturePerfectPassiveIndicative = firstFuturePerfectPassiveIndicative

var thirdImperfectPassiveSubjunctive    = firstImperfectPassiveSubjunctive
var thirdPerfectPassiveSubjunctive      = firstPerfectPassiveSubjunctive
var thirdPluperfectPassiveSubjunctive   = firstPluperfectPassiveSubjunctive

var ThirdConjugationImperatives = []ImperativePattern{

	// Present Active
	{"e", 2, "singular", "present", "active"},
	{"ite", 2, "plural", "present", "active"},

	// Future Active
	{"itō", 2, "singular", "future", "active"},
	{"itōte", 2, "plural", "future", "active"},
	{"itō", 3, "singular", "future", "active"},
	{"untō", 3, "plural", "future", "active"},

	// Present Passive
	{"ere", 2, "singular", "present", "passive"},
	{"iminī", 2, "plural", "present", "passive"},

	// Future Passive
	{"itor", 2, "singular", "future", "passive"},
	{"iminī", 2, "plural", "future", "passive"},
	{"itor", 3, "singular", "future", "passive"},
	{"untor", 3, "plural", "future", "passive"},
}

var ThirdConjugationPatterns = []FinitePattern{

	// ACTIVE INDICATIVE

	{
		Stem:    "present",
		Endings: thirdPresentActiveIndicative,
		Tense:   "present",
		Mood:    "indicative",
		Voice:   "active",
	},
	{
		Stem:    "present",
		Endings: thirdImperfectActiveIndicative,
		Tense:   "imperfect",
		Mood:    "indicative",
		Voice:   "active",
	},
	{
		Stem:    "present",
		Endings: thirdFutureActiveIndicative,
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
		Endings: thirdPresentActiveSubjunctive,
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
		Endings: thirdPresentPassiveIndicative,
		Tense:   "present",
		Mood:    "indicative",
		Voice:   "passive",
	},
	{
		Stem:    "present",
		Endings: thirdImperfectPassiveIndicative,
		Tense:   "imperfect",
		Mood:    "indicative",
		Voice:   "passive",
	},
	{
		Stem:    "present",
		Endings: thirdFuturePassiveIndicative,
		Tense:   "future",
		Mood:    "indicative",
		Voice:   "passive",
	},

	// PASSIVE SUBJUNCTIVE

	{
		Stem:    "present",
		Endings: thirdPresentPassiveSubjunctive,
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

var ThirdConjugationPerfectPassivePatterns = []PerfectPassivePattern{
	{
		Endings: thirdPerfectPassiveIndicative,
		Tense:   "perfect",
		Mood:    "indicative",
	},
	{
		Endings: thirdPluperfectPassiveIndicative,
		Tense:   "pluperfect",
		Mood:    "indicative",
	},
	{
		Endings: thirdFuturePerfectPassiveIndicative,
		Tense:   "future perfect",
		Mood:    "indicative",
	},
	{
		Endings: thirdPerfectPassiveSubjunctive,
		Tense:   "perfect",
		Mood:    "subjunctive",
	},
	{
		Endings: thirdPluperfectPassiveSubjunctive,
		Tense:   "pluperfect",
		Mood:    "subjunctive",
	},
}