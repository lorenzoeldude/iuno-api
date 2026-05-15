package dictionary

import (
	"database/sql"
	"fmt"
)

type Dictionary struct {
	DB *sql.DB
}

type Word struct {
	ID         int
	Lemma      string
	Declension sql.NullInt64
	Meanings   []Meaning
}

type Meaning struct {
	ID        int
	WordID    int
	Translation string
}

// -------------------------
// GET WORD BY LEMMA
// -------------------------

func (d *Dictionary) GetWordByLemma(lemma string) (*Word, error) {
	query := `
		SELECT id, lemma, declension
		FROM words
		WHERE lemma = $1
	`

	var w Word

	err := d.DB.QueryRow(query, lemma).Scan(
		&w.ID,
		&w.Lemma,
		&w.Declension, // ✅ handles NULL safely
	)
	if err != nil {
		return nil, err
	}

	meanings, err := d.getMeaningsByWordID(w.ID)
	if err != nil {
		return nil, err
	}

	w.Meanings = meanings

	return &w, nil
}

// -------------------------
// GET MULTIPLE WORDS (SEARCH)
// -------------------------

func (d *Dictionary) SearchWords(prefix string) ([]Word, error) {
	query := `
		SELECT id, lemma, declension
		FROM words
		WHERE lemma ILIKE $1
		LIMIT 20
	`

	rows, err := d.DB.Query(query, prefix+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []Word

	for rows.Next() {
		var w Word

		err := rows.Scan(
			&w.ID,
			&w.Lemma,
			&w.Declension,
		)
		if err != nil {
			return nil, err
		}

		meanings, err := d.getMeaningsByWordID(w.ID)
		if err != nil {
			return nil, err
		}

		w.Meanings = meanings
		words = append(words, w)
	}

	return words, nil
}

// -------------------------
// GET MEANINGS
// -------------------------

func (d *Dictionary) getMeaningsByWordID(wordID int) ([]Meaning, error) {
	query := `
		SELECT id, word_id, translation
		FROM meanings
		WHERE word_id = $1
	`

	rows, err := d.DB.Query(query, wordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meanings []Meaning

	for rows.Next() {
		var m Meaning

		err := rows.Scan(
			&m.ID,
			&m.WordID,
			&m.Translation,
		)
		if err != nil {
			return nil, err
		}

		meanings = append(meanings, m)
	}

	return meanings, nil
}

// -------------------------
// OPTIONAL: CREATE WORD
// -------------------------

func (d *Dictionary) CreateWord(lemma string, declension *int) (int, error) {
	query := `
		INSERT INTO words (lemma, declension)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int

	err := d.DB.QueryRow(query, lemma, declension).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert word failed: %w", err)
	}

	return id, nil
}

// -------------------------
// OPTIONAL: ADD MEANING
// -------------------------

func (d *Dictionary) AddMeaning(wordID int, translation string) error {
	query := `
		INSERT INTO meanings (word_id, translation)
		VALUES ($1, $2)
	`

	_, err := d.DB.Exec(query, wordID, translation)
	if err != nil {
		return fmt.Errorf("insert meaning failed: %w", err)
	}

	return nil
}