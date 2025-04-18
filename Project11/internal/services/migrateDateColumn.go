package services

import (
	"fmt"
	"log"
)

func (s *Storage) MigrateDateColumn(data map[string]interface{}) (any, error) {
	log.Println("Migrating Date Column")

	rows, err := s.QueryWithTiming(`
		SELECT sch_metadata_id as type_id, 'metadata' as word_type, 
		       upper(column_alias) as word, 'aDD' AS category
		FROM sch_metadata sm 
		JOIN sch_column_date_format cdf ON sm.sch_schema_id = cdf.sch_schema_id 
			AND sch_metadata_id = metadata_id
		WHERE sm.sch_schema_id = ? AND column_is_date = 1 AND column_is_active = 1 
		ORDER BY is_default_date DESC 
		LIMIT 1;`, data["domainId"])

	if err != nil {
		return nil, fmt.Errorf("error fetching default date format: %w", err)
	}

	_addRow, err := GetFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("error processing query results from first query: %w", err)
	}

	if len(_addRow) == 0 {
		rows, err = s.QueryWithTiming(`
			SELECT sch_metadata_id
			FROM sch_metadata
			WHERE column_is_date = 1 AND sch_schema_id = ? AND column_name IS NOT NULL
			LIMIT 1;`, data["domainId"])

		if err != nil {
			return nil, fmt.Errorf("error fetching metadata_id: %w", err)
		}

		_newRow, err := GetFromRows(rows)
		if err != nil || len(_newRow) == 0 {
			return nil, fmt.Errorf("error processing second query results or no rows found")
		}

		schMetadataID, ok := _newRow[0]["sch_metadata_id"].(int)
		if !ok {
			return nil, fmt.Errorf("unable to extract sch_metadata_id")
		}

		sqlClearDefaultDate := `UPDATE sch_column_date_format SET is_default_date = 0 WHERE sch_schema_id = ?`
		_, err = s.Db.Exec(sqlClearDefaultDate, data["domainId"])
		if err != nil {
			log.Fatal("Error executing update:", err)
		}

		_, err = s.Db.Exec(`
			INSERT INTO sch_column_date_format(metadata_id, sch_schema_id, is_default_date, modified_by, is_active)
			VALUES (?, ?, ?, ?, 1)
			ON DUPLICATE KEY UPDATE metadata_id = ?, is_default_date = ?`,
			schMetadataID, data["domainId"], 1, data["userId"], schMetadataID, 1)

		if err != nil {
			return nil, fmt.Errorf("error executing insert/update query: %w", err)
		}

		rows, err = s.QueryWithTiming(`
			SELECT sch_metadata_id as type_id, 'metadata' as word_type,
			       upper(column_alias) as word, 'aDD' AS category
			FROM sch_metadata sm
			JOIN sch_column_date_format cdf ON sm.sch_schema_id = cdf.sch_schema_id
				AND sch_metadata_id = metadata_id
			WHERE sm.sch_schema_id = ? AND column_is_date = 1 AND column_is_active = 1
			ORDER BY is_default_date DESC
			LIMIT 1;`, data["domainId"])

		if err != nil {
			return nil, fmt.Errorf("error re-fetching default date format: %w", err)
		}

		_addRow, err = GetFromRows(rows)
		if err != nil || len(_addRow) == 0 {
			return nil, fmt.Errorf("error processing re-fetched query results or no rows found")
		}
	}
	go s.hPublish(fmt.Sprintf("%0.f", data["domainId"]), "aDD", _addRow[0])

	return _addRow[0], nil
}
