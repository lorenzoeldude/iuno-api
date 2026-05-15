package dictionary

import (
	"database/sql"
)

func GetMeanings(db *sql.DB, lemmaID int) ([]string, error) {

	rows, err := db.Query(`
		SELECT meaning
		FROM meanings
		WHERE lemma_id = $1
		ORDER BY sort_order ASC
	`, lemmaID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meanings []string

	for rows.Next() {

		var m string

		if err := rows.Scan(&m); err != nil {
			return nil, err
		}

		meanings = append(meanings, m)
	}

	return meanings, nil
}