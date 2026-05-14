package morphology

import "iuno-api/models"

import "strings"

//
// VERB STEM SYSTEM
//
// Latin verbs are generated from stems.
//
// Example:
//
// amare
//
// Present stem: ama-
// Perfect stem: amav-
// Supine stem: amat-
//
// This file is responsible ONLY
// for building those stems.
//

func buildVerbStems(word models.Word) VerbStems {

	// -------------------------
	// USE DATABASE STEMS IF THEY EXIST
	// -------------------------

	present := word.Stem
	perfect := word.Perfect
	supine := word.Supine

	// -------------------------
	// FALLBACK STEM GENERATION
	// -------------------------

	// PRESENT STEM
	//
	// 1st conjugation:
	// amare -> ama-
	//

	if present == "" {

		switch word.Conjugation {

		case 1:

			if strings.HasSuffix(word.Lemma, "are") {
				present = strings.TrimSuffix(
					word.Lemma,
					"are",
				)
			}

		case 2:

			if strings.HasSuffix(word.Lemma, "ere") {
				present = strings.TrimSuffix(
					word.Lemma,
					"ere",
				)
			}

		case 3:

			if strings.HasSuffix(word.Lemma, "ere") {
				present = strings.TrimSuffix(
					word.Lemma,
					"ere",
				)
			}

		case 4:

			if strings.HasSuffix(word.Lemma, "ire") {
				present = strings.TrimSuffix(
					word.Lemma,
					"ire",
				)
			}
		}
	}

	// -------------------------
	// PERFECT STEM
	//
	// amavi -> amav-
	// monui -> monu-
	//

	if perfect == "" {

		if strings.HasSuffix(word.Lemma, "avi") {
			perfect = strings.TrimSuffix(
				word.Lemma,
				"i",
			)
		}

		if strings.HasSuffix(word.Lemma, "ui") {
			perfect = strings.TrimSuffix(
				word.Lemma,
				"i",
			)
		}
	}

	// -------------------------
	// SUPINE STEM
	//
	// amatum -> amat-
	//

	if supine == "" {

		if strings.HasSuffix(word.Lemma, "tum") {
			supine = strings.TrimSuffix(
				word.Lemma,
				"um",
			)
		}
	}

	// -------------------------
	// RETURN STEM SYSTEM
	// -------------------------

	return VerbStems{
		Present: present,
		Perfect: perfect,
		Supine:  supine,
	}
}