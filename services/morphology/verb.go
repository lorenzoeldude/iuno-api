package morphology

// import "iuno-api/models"

// // VERB ENGINE
// // It does NOT contain the actual
// // conjugation rules themselves.
// //
// // Responsibilities:
// // 1. Build stem system
// // 2. Generate finite forms
// // 3. Generate non-finite forms
// // 4. Merge results

// func GenerateVerb(lemma models.Lemma) []models.Form {

// 	// BUILD STEM SYSTEM
// 	stems := buildVerbStems(lemma)
// 	var forms []models.Form

// 	// FINITE FORMS
// 	// indicative
// 	// subjunctive
// 	// active
// 	// passive
// 	// all tenses/persons
// 	finiteForms := generateFiniteForms(
// 		lemma,
// 		stems,
// 	)

// 	forms = append(forms, finiteForms...)

// 	// NON-FINITE FORMS
// 	// infinitives
// 	// participles
// 	// gerunds
// 	// gerundives
// 	// supines
// 	nonFiniteForms := generateNonFiniteForms(
// 		lemma,
// 		stems,
// 	)

// 	forms = append(forms, nonFiniteForms...)

// 	return forms
// }