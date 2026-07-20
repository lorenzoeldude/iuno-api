package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"math/rand"

	"iuno-api/db"
)

type LessonVocabularyRequest struct {
	Vocabulary []string `json:"vocabulary"`
}

func UpdateLessonVocabularyHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("=== UpdateLessonVocabularyHandler ===")
	log.Println("Method:", r.Method)
	log.Println("Path:", r.URL.Path)

	if r.Method != http.MethodPut {
		log.Println("ERROR: wrong method")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idString := strings.TrimPrefix(
		r.URL.Path,
		"/api/admin/lessons/",
	)

	idString = strings.TrimSuffix(
		idString,
		"/vocabulary",
	)

	log.Println("Parsed lesson id string:", idString)

	lessonID, err := strconv.Atoi(idString)
	if err != nil {
		log.Println("ERROR: invalid lesson id:", err)
		http.Error(
			w,
			"invalid lesson id",
			http.StatusBadRequest,
		)
		return
	}

	log.Println("Lesson ID:", lessonID)

	var req LessonVocabularyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("ERROR: decode failed:", err)
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	log.Println("Vocabulary received:", req.Vocabulary)

	ctx := context.Background()

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		log.Println("ERROR: begin transaction:", err)
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		`DELETE FROM lesson_vocabulary WHERE lesson_id = $1`,
		lessonID,
	)

	if err != nil {
		log.Println("ERROR: delete old vocabulary:", err)
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	log.Println("Deleted existing vocabulary")

	for _, lemma := range req.Vocabulary {

		lemma = strings.TrimSpace(lemma)

		if lemma == "" {
			continue
		}

		log.Println("Looking up lemma:", lemma)

		var lemmaID int

		err = tx.QueryRow(
			ctx,
			`SELECT id FROM lemmas WHERE lemma = $1`,
			lemma,
		).Scan(&lemmaID)

		if err != nil {
			log.Println("ERROR: lemma not found:", lemma, err)
			http.Error(
				w,
				"lemma not found: "+lemma,
				http.StatusBadRequest,
			)
			return
		}

		log.Printf("Found lemma %q with id %d\n", lemma, lemmaID)

		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO lesson_vocabulary (
				lesson_id,
				lemma_id
			)
			VALUES ($1, $2)
			`,
			lessonID,
			lemmaID,
		)

		if err != nil {
			log.Println("ERROR: insert failed:", err)
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		log.Println("Inserted:", lemma)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("ERROR: commit failed:", err)
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	log.Println("Vocabulary updated successfully")

	w.WriteHeader(http.StatusNoContent)
}

func GetLessonVocabularyHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	idString := strings.TrimPrefix(
		r.URL.Path,
		"/api/admin/lessons/",
	)

	idString = strings.TrimSuffix(
		idString,
		"/vocabulary",
	)

	lessonID, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(
			w,
			"invalid lesson id",
			http.StatusBadRequest,
		)
		return
	}

	rows, err := db.Pool.Query(
		context.Background(),
		`
		SELECT l.lemma
		FROM lesson_vocabulary lv
		JOIN lemmas l ON l.id = lv.lemma_id
		WHERE lv.lesson_id = $1
		ORDER BY l.lemma
		`,
		lessonID,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	defer rows.Close()

	vocabulary := []string{}

	for rows.Next() {
		var lemma string

		if err := rows.Scan(&lemma); err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		vocabulary = append(vocabulary, lemma)
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(map[string][]string{
		"vocabulary": vocabulary,
	})
}

type LessonVocabularyQuestion struct {
	Lemma            string   `json:"lemma"`
	LemmaNormalized  string   `json:"lemma_normalized"`
	Infinitive       *string  `json:"infinitive"`
	Correct          string   `json:"correct"`
	Answers          []string `json:"answers"`
}

func GetLessonVocabularyTrainerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rand.Seed(time.Now().UnixNano())

	idString := strings.TrimPrefix(r.URL.Path, "/api/lessons/")
	idString = strings.TrimSuffix(idString, "/vocabulary")

	lessonID, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalid lesson id", http.StatusBadRequest)
		return
	}

	rows, err := db.Pool.Query(context.Background(), `
		SELECT
			l.id,
			l.lemma,
			l.lemma_normalized,
			l.infinitive,
			l.part_of_speech
		FROM lesson_vocabulary lv
		JOIN lemmas l
			ON l.id = lv.lemma_id
		WHERE lv.lesson_id = $1
	`, lessonID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	questions := []LessonVocabularyQuestion{}

	for rows.Next() {

		var lemmaID int
		var lemma string
		var lemmaNormalized string
		var infinitive *string
		var partOfSpeech string

		if err := rows.Scan(
			&lemmaID,
			&lemma,
			&lemmaNormalized,
			&infinitive,
			&partOfSpeech,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var correct string

		err = db.Pool.QueryRow(context.Background(), `
			SELECT meaning
			FROM meanings
			WHERE lemma_id = $1
			ORDER BY id
			LIMIT 1
		`, lemmaID).Scan(&correct)

		if err != nil {
			continue
		}

		wrongRows, err := db.Pool.Query(context.Background(), `
			SELECT m.meaning
			FROM meanings m
			JOIN lemmas l
				ON l.id = m.lemma_id
			WHERE
				m.lemma_id != $1
				AND l.part_of_speech = $2
			ORDER BY RANDOM()
			LIMIT 20
		`, lemmaID, partOfSpeech)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		answers := []string{correct}

		used := map[string]bool{
			correct: true,
		}

		for wrongRows.Next() {

			var wrong string

			if err := wrongRows.Scan(&wrong); err != nil {
				continue
			}

			if used[wrong] {
				continue
			}

			used[wrong] = true
			answers = append(answers, wrong)

			if len(answers) == 4 {
				break
			}
		}

		wrongRows.Close()

		if len(answers) < 4 {
			continue
		}

		rand.Shuffle(len(answers), func(i, j int) {
			answers[i], answers[j] = answers[j], answers[i]
		})

		questions = append(questions, LessonVocabularyQuestion{
			Lemma:           lemma,
			LemmaNormalized: lemmaNormalized,
			Infinitive:      infinitive,
			Correct:         correct,
			Answers:         answers,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

func LessonTrainerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	idString := strings.TrimPrefix(
		r.URL.Path,
		"/api/lessons/",
	)

	idString = strings.TrimSuffix(
		idString,
		"/trainer/random",
	)

	lessonID, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(
			w,
			"invalid lesson id",
			http.StatusBadRequest,
		)
		return
	}

	var lemmaID int
	var lemma string
	var lemmaNormalized string
	var infinitive *string
	var partOfSpeech string

	err = db.Pool.QueryRow(
		r.Context(),
		`
		SELECT
			l.id,
			l.lemma,
			l.lemma_normalized,
			l.infinitive,
			l.part_of_speech
		FROM lesson_vocabulary lv
		JOIN lemmas l
			ON l.id = lv.lemma_id
		WHERE lv.lesson_id = $1
		ORDER BY RANDOM()
		LIMIT 1
		`,
		lessonID,
	).Scan(
		&lemmaID,
		&lemma,
		&lemmaNormalized,
		&infinitive,
		&partOfSpeech,
	)

	if err != nil {
		http.Error(
			w,
			"no vocabulary found",
			http.StatusBadRequest,
		)
		return
	}

	question, err := BuildTrainerQuestion(
		r.Context(),
		lemmaID,
		lemma,
		lemmaNormalized,
		infinitive,
		partOfSpeech,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(question)
}