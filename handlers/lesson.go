package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"iuno-api/db"
	"iuno-api/models"
)


func CreateLessonHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var lesson models.Lesson

	if err := json.NewDecoder(r.Body).Decode(&lesson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.Pool.QueryRow(
		context.Background(),
		`
		INSERT INTO lessons (
			title,
			introduction,
			text,
			grammar,
			exam,
			is_published,
			image
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
		`,
		lesson.Title,
		lesson.Introduction,
		lesson.Text,
		lesson.Grammar,
		lesson.Exam,
		lesson.IsPublished,
		lesson.Image,
	).Scan(
		&lesson.ID,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(lesson)
}



func GetLessonHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}


	idString := strings.TrimPrefix(
		r.URL.Path,
		"/api/admin/lessons/",
	)

	idString = strings.TrimPrefix(
		idString,
		"/api/lessons/",
	)


	id, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(
			w,
			"invalid lesson id",
			http.StatusBadRequest,
		)
		return
	}


	var lesson models.Lesson


	err = db.Pool.QueryRow(
		context.Background(),
		`
		SELECT
			id,
			title,
			introduction,
			text,
			grammar,
			exam,
			is_published,
			image,
			created_at,
			updated_at
		FROM lessons
		WHERE id = $1
		`,
		id,
	).Scan(
		&lesson.ID,
		&lesson.Title,
		&lesson.Introduction,
		&lesson.Text,
		&lesson.Grammar,
		&lesson.Exam,
		&lesson.IsPublished,
		&lesson.Image,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)


	if err != nil {
		http.Error(
			w,
			"lesson not found",
			http.StatusNotFound,
		)
		return
	}


	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(lesson)
}

func UpdateLessonHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}


	idString := strings.TrimPrefix(
		r.URL.Path,
		"/api/admin/lessons/",
	)


	id, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(
			w,
			"invalid lesson id",
			http.StatusBadRequest,
		)
		return
	}


	var lesson models.Lesson

	if err := json.NewDecoder(r.Body).Decode(&lesson); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}


	err = db.Pool.QueryRow(
		context.Background(),
		`
		UPDATE lessons
		SET
			title = $1,
			introduction = $2,
			text = $3,
			grammar = $4,
			exam = $5,
			is_published = $6,
			image = $7,
			updated_at = NOW()
		WHERE id = $8
		RETURNING
			id,
			created_at,
			updated_at
		`,
		lesson.Title,
		lesson.Introduction,
		lesson.Text,
		lesson.Grammar,
		lesson.Exam,
		lesson.IsPublished,
		lesson.Image,
		id,
	).Scan(
		&lesson.ID,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
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

	json.NewEncoder(w).Encode(lesson)
}

func GetLessonsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Pool.Query(
		context.Background(),
		`
		SELECT
			id,
			title,
			is_published,
			image
		FROM lessons
		ORDER BY id
		`,
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

	var lessons []models.Lesson

	for rows.Next() {

		var lesson models.Lesson

		err := rows.Scan(
			&lesson.ID,
			&lesson.Title,
			&lesson.IsPublished,
			&lesson.Image,
		)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		lessons = append(lessons, lesson)
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(lessons)
}