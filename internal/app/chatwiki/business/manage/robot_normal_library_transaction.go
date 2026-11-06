// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"database/sql"
	"fmt"
	"github.com/zhimaAi/go_tools/msql"
	"sort"
	"strings"
)

type normalLibraryResult struct {
	LibraryID int  `json:"library_id"`
	Created   bool `json:"created"`
	robotKey  string
	data      msql.Datas
}

const normalLibraryEnsureMaxAttempts = 3

// Avoid holding a robot row lock with SELECT ... FOR UPDATE: read without a lock,
// then claim the binding with a conditional atomic UPDATE. On conflict, roll back
// the transaction, including any new library, and read again.
func ensureNormalLibraryTx(ctx context.Context, db *sql.DB, adminID, robotID int, prepare func(string) (msql.Datas, error)) (normalLibraryResult, error) {
	var lastErr error
	for attempt := 0; attempt < normalLibraryEnsureMaxAttempts; attempt++ {
		result, retry, err := ensureNormalLibraryOnce(ctx, db, adminID, robotID, prepare)
		if err != nil {
			return normalLibraryResult{}, err
		}
		if !retry {
			return result, nil
		}
		lastErr = fmt.Errorf("normal library bind conflict")
	}
	return normalLibraryResult{}, lastErr
}

func ensureNormalLibraryOnce(ctx context.Context, db *sql.DB, adminID, robotID int, prepare func(string) (msql.Datas, error)) (normalLibraryResult, bool, error) {
	var result normalLibraryResult
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return result, false, err
	}
	defer tx.Rollback()

	var name string
	var applicationType, defaultID, normalID int
	err = tx.QueryRowContext(ctx, `SELECT robot_key,robot_name,application_type,COALESCE(default_library_id,0),COALESCE(normal_library_id,0)
 FROM chat_ai_robot WHERE id=$1 AND admin_user_id=$2`, robotID, adminID).
		Scan(&result.robotKey, &name, &applicationType, &defaultID, &normalID)
	if err != nil {
		return result, false, err
	}
	if applicationType != define.ApplicationTypeChat {
		return result, false, sql.ErrNoRows
	}

	validLibrary := func(id int) (bool, error) {
		if id <= 0 {
			return false, nil
		}
		var exists bool
		err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM chat_ai_library WHERE id=$1 AND admin_user_id=$2 AND type=0)`, id, adminID).Scan(&exists)
		return exists, err
	}

	valid, err := validLibrary(normalID)
	if err != nil {
		return result, false, err
	}
	if !valid {
		normalID = 0
		valid, err = validLibrary(defaultID)
		if err != nil {
			return result, false, err
		}
		if valid {
			normalID = defaultID
		}
	}

	var data msql.Datas
	created := false
	if normalID == 0 {
		data, err = prepare(name)
		if err != nil {
			return result, false, err
		}
		keys := make([]string, 0, len(data))
		for key := range data {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		columns, seats, args := make([]string, 0, len(keys)), make([]string, 0, len(keys)), make([]any, 0, len(keys))
		for index, key := range keys {
			columns = append(columns, `"`+key+`"`)
			seats = append(seats, fmt.Sprintf("$%d", index+1))
			args = append(args, data[key])
		}
		err = tx.QueryRowContext(ctx, `INSERT INTO chat_ai_library (`+strings.Join(columns, ",")+`) VALUES (`+strings.Join(seats, ",")+`) RETURNING id`, args...).Scan(&normalID)
		if err != nil {
			return result, false, err
		}
		created = true
	}

	// Update only if the candidate exists and normal_library_id has no other valid binding.
	// The atomic UPDATE lets only one concurrent request claim the binding; losers roll back and read again.
	updateResult, err := tx.ExecContext(ctx, `UPDATE chat_ai_robot
SET normal_library_id=$1,
    library_ids=(
        SELECT string_agg(value, ',' ORDER BY value)
        FROM (
            SELECT DISTINCT value
            FROM unnest(string_to_array(COALESCE(library_ids, ''), ',') || ARRAY[$1::text]) AS value
            WHERE value <> ''
        ) merged
    )
WHERE id=$2 AND admin_user_id=$3
  AND EXISTS (
      SELECT 1 FROM chat_ai_library candidate
      WHERE candidate.id=$1 AND candidate.admin_user_id=$3 AND candidate.type=0
  )
  AND (
      COALESCE(normal_library_id, 0)=0
      OR normal_library_id=$1
      OR NOT EXISTS (
          SELECT 1 FROM chat_ai_library current
          WHERE current.id=chat_ai_robot.normal_library_id
            AND current.admin_user_id=chat_ai_robot.admin_user_id
            AND current.type=0
      )
  )`, normalID, robotID, adminID)
	if err != nil {
		return result, false, err
	}
	affected, err := updateResult.RowsAffected()
	if err != nil {
		return result, false, err
	}
	if affected == 0 {
		return normalLibraryResult{}, true, nil
	}
	if err = tx.Commit(); err != nil {
		return result, false, err
	}
	result.LibraryID = normalID
	result.Created = created
	result.data = data
	return result, false, nil
}
