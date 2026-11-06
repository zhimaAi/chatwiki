// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/zhimaAi/go_tools/msql"
)

type bm25IndexMatch struct {
	IndexID    string
	DataID     string
	Similarity float64
}

// normalizeBM25Score maps positive finite BM25 scores to (0, 1) using score / (score + 1).
// It returns zero for invalid scores and clamps rounding to keep the result below one.
func normalizeBM25Score(score float64) float64 {
	if score <= 0 || math.IsNaN(score) || math.IsInf(score, 0) {
		return 0
	}
	similarity := score / (score + 1)
	if similarity >= 1 {
		return math.Nextafter(1, 0)
	}
	return similarity
}

// getBM25IndexMatches retrieves BM25-ranked index rows from the selected libraries using bound parameters.
// It excludes deleted index rows and filters normalized scores by the similarity threshold.
// Returned matches retain the index and paragraph IDs for subsequent paragraph retrieval.
func getBM25IndexMatches(question, libraryIDs string, size int, similarity float64) ([]bm25IndexMatch, error) {
	result := make([]bm25IndexMatch, 0)
	question = strings.TrimSpace(question)
	if question == `` || strings.TrimSpace(libraryIDs) == `` || size <= 0 {
		return result, nil
	}
	if math.IsNaN(similarity) || similarity < 0 || similarity > 1 {
		return nil, fmt.Errorf("invalid BM25 similarity threshold: %v", similarity)
	}

	args := []any{question}
	placeholders := make([]string, 0)
	for _, value := range strings.Split(libraryIDs, `,`) {
		id, err := strconv.ParseInt(strings.TrimSpace(value), 10, 32)
		if err != nil || id <= 0 {
			return nil, fmt.Errorf("invalid BM25 library ID: %q", value)
		}
		args = append(args, id)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
	}
	args = append(args, size)
	query := fmt.Sprintf(`SELECT id, data_id, pdb.score(id)::double precision AS bm25_score
		FROM chat_ai_library_file_data_index
		WHERE content ||| $1
		  AND library_id IN (%s)
		  AND delete_time = 0
		ORDER BY pdb.score(id) DESC NULLS LAST, id ASC
		LIMIT $%d`, strings.Join(placeholders, `,`), len(args))

	db, err := msql.GetDB(define.Postgres)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("BM25 recall query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var indexID, dataID int64
		var score sql.NullFloat64
		if err := rows.Scan(&indexID, &dataID, &score); err != nil {
			return nil, fmt.Errorf("scan BM25 recall result: %w", err)
		}
		if !score.Valid || indexID <= 0 || dataID <= 0 {
			continue
		}
		normalized := normalizeBM25Score(score.Float64)
		if normalized <= 0 || normalized < similarity {
			continue
		}
		result = append(result, bm25IndexMatch{
			IndexID:    strconv.FormatInt(indexID, 10),
			DataID:     strconv.FormatInt(dataID, 10),
			Similarity: normalized,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read BM25 recall results: %w", err)
	}
	return result, nil
}
