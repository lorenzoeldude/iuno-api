package handlers

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"iuno-api/db"
	"iuno-api/middleware"
	"iuno-api/utils"
)

func ListTrainerHandler(w http.ResponseWriter, r *http.Request) {


	if r.Method != http.MethodGet {

		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)

		return
	}



	claimsRaw := r.Context().Value(middleware.UserContextKey)

	claims, ok := claimsRaw.(*utils.Claims)

	if !ok || claims == nil {

		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)

		return
	}



	userID := claims.UserID



	// =====================================================
	// 1. PICK RANDOM LEMMA FROM USER LIST
	// =====================================================

	var lemmaID int
	var lemma string
	var lemmaNormalized string
	var infinitive *string
	var partOfSpeech string


	err := db.Pool.QueryRow(r.Context(), `
		SELECT
			l.id,
			l.lemma,
			l.lemma_normalized,
			l.infinitive,
			l.part_of_speech
		FROM word_list_lemmas wll
		JOIN word_lists wl 
			ON wl.id = wll.list_id
		JOIN lemmas l 
			ON l.id = wll.lemma_id
		WHERE wl.user_id = $1
		ORDER BY RANDOM()
		LIMIT 1
	`,
		userID,
	).Scan(
		&lemmaID,
		&lemma,
		&lemmaNormalized,
		&infinitive,
		&partOfSpeech,
	)


	if err != nil {

		log.Println("NO WORDS ERROR:", err)

		http.Error(
			w,
			"no words in list",
			http.StatusBadRequest,
		)

		return
	}





	// =====================================================
	// 2. GET CORRECT MEANING
	// =====================================================

	var correct string


	err = db.Pool.QueryRow(r.Context(), `
		SELECT meaning
		FROM meanings
		WHERE lemma_id = $1
		ORDER BY sort_order ASC
		LIMIT 1
	`,
		lemmaID,
	).Scan(&correct)



	if err != nil {

		log.Println("MEANING FETCH ERROR:", err)

		http.Error(
			w,
			"no meaning found",
			http.StatusBadRequest,
		)

		return
	}





	// =====================================================
	// 3. GET DEFINITION
	// =====================================================

	var definition string


	err = db.Pool.QueryRow(r.Context(), `
		SELECT definition
		FROM definitions
		WHERE lemma_id = $1
		ORDER BY id
		LIMIT 1
	`,
		lemmaID,
	).Scan(&definition)


	if err != nil {

		definition = ""

	}





	// =====================================================
	// 4. GET EXAMPLES
	// =====================================================

	examples := []string{}


	exampleRows, err := db.Pool.Query(r.Context(), `
		SELECT example
		FROM examples
		WHERE lemma_id = $1
		ORDER BY id
		LIMIT 3
	`,
		lemmaID,
	)



	if err == nil {

		defer exampleRows.Close()


		for exampleRows.Next() {

			var example string


			if err := exampleRows.Scan(&example); err == nil {

				examples = append(
					examples,
					example,
				)

			}

		}

	}





	// =====================================================
	// 5. GET WRONG ANSWERS
	// =====================================================

	rows, err := db.Pool.Query(r.Context(), `
		SELECT m.meaning
		FROM meanings m
		JOIN lemmas l ON l.id = m.lemma_id
		WHERE
			m.lemma_id != $1
			AND l.part_of_speech = $2
		ORDER BY RANDOM()
		LIMIT 20
	`,
		lemmaID,
		partOfSpeech,
	)



	if err != nil {

		log.Println("WRONG MEANINGS ERROR:", err)

		http.Error(
			w,
			"db error",
			http.StatusInternalServerError,
		)

		return
	}


	defer rows.Close()



	answers := []string{
		correct,
	}



	used := map[string]bool{
		correct: true,
	}



	for rows.Next() {


		var meaning string


		if err := rows.Scan(&meaning); err != nil {

			continue

		}



		if used[meaning] {

			continue

		}



		used[meaning] = true

		answers = append(
			answers,
			meaning,
		)



		if len(answers) == 3 {

			break

		}

	}





	// fallback

	for len(answers) < 3 {

		answers = append(
			answers,
			"—",
		)

	}





	// =====================================================
	// 6. SHUFFLE
	// =====================================================

	rand.Shuffle(
		len(answers),
		func(i, j int) {

			answers[i], answers[j] =
				answers[j], answers[i]

		},
	)





	// =====================================================
	// RESPONSE
	// =====================================================

	w.Header().Set(
		"Content-Type",
		"application/json",
	)


	json.NewEncoder(w).Encode(
		TrainerQuestion{

			Lemma: lemma,

			LemmaNormalized: lemmaNormalized,

			Infinitive: infinitive,

			Correct: correct,

			Answers: answers,

			Definition: definition,

			Examples: examples,

		},
	)

}