package handlers

import (
	"encoding/xml"
	"net/http"

	"iuno-api/db"
)

type SitemapURL struct {
	Loc string `xml:"loc"`
}

type SitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

func SitemapHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	rows, err := db.Pool.Query(
		r.Context(),
		`
		SELECT lemma_normalized
		FROM lemmas
		WHERE lemma_normalized IS NOT NULL
		  AND lemma_normalized != ''
		ORDER BY lemma_normalized
		`,
	)

	if err != nil {
		http.Error(
			w,
			"failed to query dictionary",
			http.StatusInternalServerError,
		)
		return
	}

	defer rows.Close()

	sitemap := SitemapURLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]SitemapURL, 0),
	}

	for rows.Next() {

		var lemma string

		if err := rows.Scan(&lemma); err != nil {
			http.Error(
				w,
				"failed to read dictionary",
				http.StatusInternalServerError,
			)
			return
		}

		sitemap.URLs = append(
			sitemap.URLs,
			SitemapURL{
				Loc: "https://www.iunoni.com/dictionary/" + lemma,
			},
		)
	}

	if err := rows.Err(); err != nil {
		http.Error(
			w,
			"failed to read dictionary",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/xml; charset=utf-8",
	)

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(
		[]byte(xml.Header),
	)

	if err != nil {
		return
	}

	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")

	if err := encoder.Encode(sitemap); err != nil {
		return
	}
}
