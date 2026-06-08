package morphology

//
// DIFFERENT FROM 1ST CONJUGATION
//

var secondPresentActiveIndicative = VerbEndings{
	"singular": {
		"first":  "eō",
		"second": "ēs",
		"third":  "et",
	},
	"plural": {
		"first":  "ēmus",
		"second": "ētis",
		"third":  "ent",
	},
}

var secondImperfectActiveIndicative = VerbEndings{
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

var secondFutureActiveIndicative = VerbEndings{
	"singular": {
		"first":  "ēbō",
		"second": "ēbis",
		"third":  "ēbit",
	},
	"plural": {
		"first":  "ēbimus",
		"second": "ēbitis",
		"third":  "ēbunt",
	},
}

var secondPresentActiveSubjunctive = VerbEndings{
	"singular": {
		"first":  "eam",
		"second": "eās",
		"third":  "eat",
	},
	"plural": {
		"first":  "eāmus",
		"second": "eātis",
		"third":  "eant",
	},
}

var secondPresentPassiveIndicative = VerbEndings{
	"singular": {
		"first":  "eor",
		"second": "ēris",
		"third":  "ētur",
	},
	"plural": {
		"first":  "ēmur",
		"second": "ēminī",
		"third":  "entur",
	},
}

var secondImperfectPassiveIndicative = VerbEndings{
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

var secondFuturePassiveIndicative = VerbEndings{
	"singular": {
		"first":  "ēbor",
		"second": "ēberis",
		"third":  "ēbitur",
	},
	"plural": {
		"first":  "ēbimur",
		"second": "ēbiminī",
		"third":  "ēbuntur",
	},
}

var secondPresentPassiveSubjunctive = VerbEndings{
	"singular": {
		"first":  "ear",
		"second": "eāris",
		"third":  "eātur",
	},
	"plural": {
		"first":  "eāmur",
		"second": "eāminī",
		"third":  "eantur",
	},
}

var secondPerfectActiveIndicative = firstPerfectActiveIndicative
var secondPluperfectActiveIndicative = firstPluperfectActiveIndicative
var secondFuturePerfectActiveIndicative = firstFuturePerfectActiveIndicative

var secondImperfectActiveSubjunctive = firstImperfectActiveSubjunctive
var secondPerfectActiveSubjunctive = firstPerfectActiveSubjunctive
var secondPluperfectActiveSubjunctive = firstPluperfectActiveSubjunctive

var secondPerfectPassiveIndicative = firstPerfectPassiveIndicative
var secondPluperfectPassiveIndicative = firstPluperfectPassiveIndicative
var secondFuturePerfectPassiveIndicative = firstFuturePerfectPassiveIndicative

var secondImperfectPassiveSubjunctive = firstImperfectPassiveSubjunctive
var secondPerfectPassiveSubjunctive = firstPerfectPassiveSubjunctive
var secondPluperfectPassiveSubjunctive = firstPluperfectPassiveSubjunctive

var SecondConjugationImperatives = []ImperativePattern{
	{"ē", 2, "singular", "present", "active"},
	{"ēte", 2, "plural", "present", "active"},

	{"ētō", 2, "singular", "future", "active"},
	{"ētōte", 2, "plural", "future", "active"},
	{"ētō", 3, "singular", "future", "active"},
	{"entō", 3, "plural", "future", "active"},

	{"ērī", 2, "singular", "present", "passive"},
	{"ēminī", 2, "plural", "present", "passive"},

	{"ētor", 2, "singular", "future", "passive"},
	{"ēminī", 2, "plural", "future", "passive"},
	{"ētor", 3, "singular", "future", "passive"},
	{"entor", 3, "plural", "future", "passive"},
}

var SecondConjugationPatterns = []FinitePattern{

	// ACTIVE INDICATIVE

	{
		Stem:    "present",
		Endings: secondPresentActiveIndicative,
		Tense:   "present",
		Mood:    "indicative",
		Voice:   "active",
	},
	{
		Stem:    "present",
		Endings: secondImperfectActiveIndicative,
		Tense:   "imperfect",
		Mood:    "indicative",
		Voice:   "active",
	},
	{
		Stem:    "present",
		Endings: secondFutureActiveIndicative,
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
		Endings: secondPresentActiveSubjunctive,
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
		Endings: secondPresentPassiveIndicative,
		Tense:   "present",
		Mood:    "indicative",
		Voice:   "passive",
	},
	{
		Stem:    "present",
		Endings: secondImperfectPassiveIndicative,
		Tense:   "imperfect",
		Mood:    "indicative",
		Voice:   "passive",
	},
	{
		Stem:    "present",
		Endings: secondFuturePassiveIndicative,
		Tense:   "future",
		Mood:    "indicative",
		Voice:   "passive",
	},

	// PASSIVE SUBJUNCTIVE

	{
		Stem:    "present",
		Endings: secondPresentPassiveSubjunctive,
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

var SecondConjugationPerfectPassivePatterns = []PerfectPassivePattern{
	{
		Endings: secondPerfectPassiveIndicative,
		Tense:   "perfect",
		Mood:    "indicative",
	},
	{
		Endings: secondPluperfectPassiveIndicative,
		Tense:   "pluperfect",
		Mood:    "indicative",
	},
	{
		Endings: secondFuturePerfectPassiveIndicative,
		Tense:   "future perfect",
		Mood:    "indicative",
	},
	{
		Endings: secondPerfectPassiveSubjunctive,
		Tense:   "perfect",
		Mood:    "subjunctive",
	},
	{
		Endings: secondPluperfectPassiveSubjunctive,
		Tense:   "pluperfect",
		Mood:    "subjunctive",
	},
}