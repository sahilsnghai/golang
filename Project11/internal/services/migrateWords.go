package services

import (
	"fmt"
	"log"
	"strings"
)

func (s *Storage) MigrateWords(data map[string]interface{}) (any, error) {
	log.Println("Migrate Words")

	domainID, ok := data["domainId"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing domainId")
	}

	rows_metadata, err := s.QueryWithTiming(`
				SELECT
					sch_metadata_id AS type_id,
					'metadata' AS word_type,
					column_alias AS word,
					CASE
						WHEN ('column_type' = 'Measure') THEN 'aM'
						WHEN ('column_type' = 'Dimension'
						AND column_is_date = 1) THEN 'aDD'
						WHEN ('column_type' = 'Dimension') THEN 'aD'
						WHEN ('column_type' = 'Calculated KPI') THEN 'aC'
						WHEN ('column_type' = 'Latitude') THEN 'aLat'
						WHEN ('column_type' = 'Longitude') THEN 'aLng'
						WHEN ('column_type' = 'Location') THEN 'aLoc'
						WHEN ('column_type' = 'vCalculated KPI') THEN 'avC'
						WHEN ('column_type' = 'vDimension') THEN 'avD'
						WHEN ('column_type' = 'vMeasure') THEN 'avM'
						WHEN ('column_type' = 'Variable') THEN 'aVar'
						WHEN ('column_type' = 'vDimension' and column_is_date = 1) THEN 'avD'
					END AS category
				FROM
					sch_metadata am
				WHERE
					am.sch_schema_id = ?
					AND
					column_is_active = 1
					AND is_active = 1
					AND am.is_global = 1;`, domainID)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata words: %w", err)
	}

	_metadata, err := GetFromRows(rows_metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to retrive rows for metadata: %w", err)
	}

	s.Publish(fmt.Sprintf("%0.f", data["domainId"]),
		fmt.Sprintf("%0.f<<-->>words<<-->>metadata", data["domainId"]),
		_metadata,
	)

	log.Println("published", fmt.Sprintf("%0.f<<-->>words<<-->>metadata", data["domainId"]))

	filterLimit := s.calculateFilterLimit()
	fmt.Printf("filterLimit %d\n", filterLimit)

	rows_filters, err := s.QueryWithTiming(`
	with metadata as (
			   SELECT
				   sch_metadata_id,
				   sch_schema_id,
				   column_alias
			   FROM
				   sch_metadata sm
			   WHERE
				   sm.enable_for_filter = 1
				   and sm.column_is_active = 1
				   and sm.is_active = 1
				   and sm.sch_schema_id = ?
			   )
	   SELECT
		   sfv.metadata_id as type_id,
		   'filter' as word_type,
		   sfv.filter_name as word,
		   m.column_alias
	   FROM
		   sch_filter_values sfv
	   INNER JOIN metadata m on
		   sfv.metadata_id = m.sch_metadata_id
		   AND sfv.sch_schema_id = m.sch_schema_id
	   LIMIT ?;`, domainID, filterLimit)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch filter words: %w", err)
	}

	_filters, err := GetFromRows(rows_filters)
	if err != nil {
		return nil, fmt.Errorf("failed to retrive rows for filters: %w", err)
	}
	s.Publish(fmt.Sprintf("%0.f", data["domainId"]),
		fmt.Sprintf("%0.f<<-->>words<<-->>filters", data["domainId"]),
		_filters,
	)

	words := make(map[string][]map[string]interface{})
	_words := append(_metadata, _filters...)

	for _, item := range _words {
		word, ok := item["word"].(string)
		if !ok {
			log.Println("Error: Missing or invalid 'word' key")
			continue
		}

		wordName := strings.ToUpper(word)

		// if _, exists := words[wordName]; exists {
		// 	words[wordName] = append(words[wordName], item)
		// } else {
		// 	words[wordName] = []map[string]interface{}{item}
		// }
		words[wordName] = append(words[wordName], item)
	}

	return words, nil
}
