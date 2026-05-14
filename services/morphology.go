package services

import (
	"context"
	"strings"

	"iuno-api/db"
)

//
// RESULT STRUCT
//
type ParseResult struct {
	Input   string `json:"input"`
	Lemma   string `json:"lemma"`
	Type    string `json:"type"`
	Grammar string `json:"grammar"`
	Meaning string `json:"meaning"`
}

//
// MAIN ENTRY POINT
//
func ParseWord(input string) ParseResult {

	word := strings.ToLower(strings.TrimSpace(input))

	// 1. irregular overrides first
	if override := checkOverride(word); override != nil {
		return *override
	}

	// 2. exact dictionary match
	if exact := matchExact(word); exact != nil {
		return *exact
	}

	// 3. rule-based morphology fallback
	return ruleEngine(word)
}

//
// STEP 4 — IRREGULAR OVERRIDES
//
func checkOverride(word string) *ParseResult {

	var r ParseResult

	err := db.Pool.QueryRow(context.Background(), `
		SELECT lemma, type, grammar, meaning
		FROM word_overrides
		WHERE form=$1
	`, word).Scan(
		&r.Lemma,
		&r.Type,
		&r.Grammar,
		&r.Meaning,
	)

	if err != nil {
		return nil
	}

	r.Input = word
	return &r
}

//
// STEP 5 — EXACT DICTIONARY MATCH
//
func matchExact(word string) *ParseResult {

	var lemma, wtype, meaning string

	err := db.Pool.QueryRow(context.Background(), `
		SELECT latin, type, translation_1
		FROM words
		WHERE latin=$1
	`, word).Scan(&lemma, &wtype, &meaning)

	if err != nil {
		return nil
	}

	return &ParseResult{
		Input:   word,
		Lemma:   lemma,
		Type:    wtype,
		Grammar: "base form",
		Meaning: meaning,
	}
}

//
// STEP 6 — RULE ENGINE (MORPHOLOGY)
//
func ruleEngine(word string) ParseResult {

	// -------------------------
	// NOUN / ADJECTIVE RULES
	// -------------------------
	switch {

	case strings.HasSuffix(word, "arum"):
		return ParseResult{
			Input:   word,
			Type:    "noun/adjective",
			Grammar: "genitive plural",
			Lemma:   guessLemma(word, 4),
		}

	case strings.HasSuffix(word, "ae"):
		return ParseResult{
			Input:   word,
			Type:    "noun/adjective",
			Grammar: "genitive singular / nominative plural",
			Lemma:   guessLemma(word, 2),
		}

	case strings.HasSuffix(word, "am"):
		return ParseResult{
			Input:   word,
			Type:    "noun/adjective",
			Grammar: "accusative singular",
			Lemma:   guessLemma(word, 2),
		}

	case strings.HasSuffix(word, "um"):
		return ParseResult{
			Input:   word,
			Type:    "noun/adjective",
			Grammar: "accusative singular (2nd declension)",
			Lemma:   guessLemma(word, 2),
		}

	// -------------------------
	// VERB RULES (basic present system)
	// -------------------------
	case strings.HasSuffix(word, "o"):
		return ParseResult{
			Input:   word,
			Type:    "verb",
			Grammar: "1st person singular present",
			Lemma:   guessLemma(word, 1),
		}

	case strings.HasSuffix(word, "t"):
		return ParseResult{
			Input:   word,
			Type:    "verb",
			Grammar: "3rd person singular present",
			Lemma:   guessLemma(word, 1),
		}

	case strings.HasSuffix(word, "nt"):
		return ParseResult{
			Input:   word,
			Type:    "verb",
			Grammar: "3rd person plural present",
			Lemma:   guessLemma(word, 2),
		}
	}

	// -------------------------
	// FALLBACK
	// -------------------------
	return ParseResult{
		Input:   word,
		Type:    "unknown",
		Grammar: "unrecognized form",
		Lemma:   word,
	}
}

//
// STEP 7 — LEMMA GUESSING
//
func guessLemma(word string, cut int) string {

	if len(word) <= cut {
		return word
	}

	return word[:len(word)-cut]
}