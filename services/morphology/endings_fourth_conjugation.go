package morphology

//
// DIFFERENT FROM 1ST/2ND CONJUGATION
//

var fourthPresentActiveIndicative = VerbEndings{
	"singular": {
		"first":  "iō",
		"second": "īs",
		"third":  "it",
	},
	"plural": {
		"first":  "īmus",
		"second": "ītis",
		"third":  "iunt",
	},
}

var fourthImperfectActiveIndicative = VerbEndings{
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

var fourthFutureActiveIndicative = VerbEndings{
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

var fourthPresentActiveSubjunctive = VerbEndings{
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

var fourthPresentPassiveIndicative = VerbEndings{
	"singular": {
		"first":  "ior",
		"second": "īris",
		"third":  "ītur",
	},
	"plural": {
		"first":  "īmur",
		"second": "īminī",
		"third":  "iuntur",
	},
}

var fourthImperfectPassiveIndicative = VerbEndings{
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

var fourthFuturePassiveIndicative = VerbEndings{
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

var fourthPresentPassiveSubjunctive = VerbEndings{
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

var fourthPerfectActiveIndicative = firstPerfectActiveIndicative
var fourthPluperfectActiveIndicative = firstPluperfectActiveIndicative
var fourthFuturePerfectActiveIndicative = firstFuturePerfectActiveIndicative

var fourthImperfectActiveSubjunctive = firstImperfectActiveSubjunctive
var fourthPerfectActiveSubjunctive = firstPerfectActiveSubjunctive
var fourthPluperfectActiveSubjunctive = firstPluperfectActiveSubjunctive

var fourthPerfectPassiveIndicative = firstPerfectPassiveIndicative
var fourthPluperfectPassiveIndicative = firstPluperfectPassiveIndicative
var fourthFuturePerfectPassiveIndicative = firstFuturePerfectPassiveIndicative

var fourthImperfectPassiveSubjunctive = firstImperfectPassiveSubjunctive
var fourthPerfectPassiveSubjunctive = firstPerfectPassiveSubjunctive
var fourthPluperfectPassiveSubjunctive = firstPluperfectPassiveSubjunctive

var FourthConjugationImperatives = []ImperativePattern{
	{"ī", 2, "singular", "present", "active"},
	{"īte", 2, "plural", "present", "active"},

	{"ītō", 2, "singular", "future", "active"},
	{"ītōte", 2, "plural", "future", "active"},
	{"ītō", 3, "singular", "future", "active"},
	{"iuntō", 3, "plural", "future", "active"},

	{"īre", 2, "singular", "present", "passive"},
	{"īminī", 2, "plural", "present", "passive"},

	{"ītor", 2, "singular", "future", "passive"},
	{"īminī", 2, "plural", "future", "passive"},
	{"ītor", 3, "singular", "future", "passive"},
	{"iuntor", 3, "plural", "future", "passive"},
}

var FourthConjugationPatterns = []FinitePattern{

	// ACTIVE INDICATIVE

	{
		Stem:    "present",
		Endings: fourthPresentActiveIndicative,
		Tense:   "present",
		Mood:    "indicative",
		Voice:   "active",
	},
	{
		Stem:    "present",
		Endings: fourthImperfectActiveIndicative,
		Tense:   "imperfect",
		Mood:    "indicative",
		Voice:   "active",
	},
	{
		Stem:    "present",
		Endings: fourthFutureActiveIndicative,
		Tense:   "future",
		Mood:    "indicative",
		Voice:   "active",
	},
	{
		Stem:    "perfect",
		Endings: fourthPerfectActiveIndicative,
		Tense:   "perfect",
		Mood:    "indicative",
		Voice:   "active",
	},
	{
		Stem:    "perfect",
		Endings: fourthPluperfectActiveIndicative,
		Tense:   "pluperfect",
		Mood:    "indicative",
		Voice:   "active",
	},
	{
		Stem:    "perfect",
		Endings: fourthFuturePerfectActiveIndicative,
		Tense:   "future perfect",
		Mood:    "indicative",
		Voice:   "active",
	},

	// ACTIVE SUBJUNCTIVE

	{
		Stem:    "present",
		Endings: fourthPresentActiveSubjunctive,
		Tense:   "present",
		Mood:    "subjunctive",
		Voice:   "active",
	},
	{
		Stem:    "imperfect_subj",
		Endings: fourthImperfectActiveSubjunctive,
		Tense:   "imperfect",
		Mood:    "subjunctive",
		Voice:   "active",
	},
	{
		Stem:    "perfect",
		Endings: fourthPerfectActiveSubjunctive,
		Tense:   "perfect",
		Mood:    "subjunctive",
		Voice:   "active",
	},
	{
		Stem:    "pluperfect_subj",
		Endings: fourthPluperfectActiveSubjunctive,
		Tense:   "pluperfect",
		Mood:    "subjunctive",
		Voice:   "active",
	},

	// PASSIVE INDICATIVE

	{
		Stem:    "present",
		Endings: fourthPresentPassiveIndicative,
		Tense:   "present",
		Mood:    "indicative",
		Voice:   "passive",
	},
	{
		Stem:    "present",
		Endings: fourthImperfectPassiveIndicative,
		Tense:   "imperfect",
		Mood:    "indicative",
		Voice:   "passive",
	},
	{
		Stem:    "present",
		Endings: fourthFuturePassiveIndicative,
		Tense:   "future",
		Mood:    "indicative",
		Voice:   "passive",
	},

	// PASSIVE SUBJUNCTIVE

	{
		Stem:    "present",
		Endings: fourthPresentPassiveSubjunctive,
		Tense:   "present",
		Mood:    "subjunctive",
		Voice:   "passive",
	},
	{
		Stem:    "imperfect_passive_subj",
		Endings: fourthImperfectPassiveSubjunctive,
		Tense:   "imperfect",
		Mood:    "subjunctive",
		Voice:   "passive",
	},
}

var FourthConjugationPerfectPassivePatterns = []PerfectPassivePattern{
	{
		Endings: fourthPerfectPassiveIndicative,
		Tense:   "perfect",
		Mood:    "indicative",
	},
	{
		Endings: fourthPluperfectPassiveIndicative,
		Tense:   "pluperfect",
		Mood:    "indicative",
	},
	{
		Endings: fourthFuturePerfectPassiveIndicative,
		Tense:   "future perfect",
		Mood:    "indicative",
	},
	{
		Endings: fourthPerfectPassiveSubjunctive,
		Tense:   "perfect",
		Mood:    "subjunctive",
	},
	{
		Endings: fourthPluperfectPassiveSubjunctive,
		Tense:   "pluperfect",
		Mood:    "subjunctive",
	},
}