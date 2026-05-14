package morphology

import "iuno-api/models"

//
// VERB ENGINE
//
// This file is the top-level orchestrator
// for ALL verb morphology.
//
// It does NOT contain the actual
// conjugation rules themselves.
//
// Responsibilities:
//
// 1. Build stem system
// 2. Generate finite forms
// 3. Generate non-finite forms
// 4. Merge results
//

func GenerateVerb(word models.Word) []models.Form {

	// -------------------------
	// BUILD STEM SYSTEM
	// -------------------------

	stems := buildVerbStems(word)

	// -------------------------
	// OUTPUT
	// -------------------------

	var forms []models.Form

	// -------------------------
	// FINITE FORMS
	//
	// indicative
	// subjunctive
	// active
	// passive
	// all tenses/persons
	// -------------------------

	finiteForms := generateFiniteForms(
		word,
		stems,
	)

	forms = append(forms, finiteForms...)

	// -------------------------
	// NON-FINITE FORMS
	//
	// infinitives
	// participles
	// gerunds
	// gerundives
	// supines
	// -------------------------

	nonFiniteForms := generateNonFiniteForms(
		word,
		stems,
	)

	forms = append(forms, nonFiniteForms...)

	// -------------------------
	// RETURN ALL GENERATED FORMS
	// -------------------------

	return forms
}