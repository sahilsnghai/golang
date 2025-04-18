package services

import (
	"fmt"
	"log"
	"sync"
)

func (s *Storage) FilterRefresh(data map[string]interface{}) (map[string]interface{}, error) {
	log.Println("Starting Filter refresh")

	// removeFilterCh := make(chan bool)
	// go func() {
	// 	removeFilterCh <- s.removeFilterValues(data)
	// }()

	metadataIDs, err := s.fetchMetadataIDs(data["domainId"])
	if err != nil || len(metadataIDs) == 0 {
		return nil, fmt.Errorf("unable to fetch metadata IDs or no filters needed: %w", err)
	}

	filterLimit := s.getFilterLimit(data)

	filters, err := fetchFilterValues(s, data["domainId"].(int), metadataIDs, filterLimit)
	if err != nil {
		return nil, fmt.Errorf("error fetching filter values: %w", err)
	}

	err = s.insertFilters(data["domainId"], data["userID"], filters)
	if err != nil {
		return nil, fmt.Errorf("error inserting filters: %w", err)
	}

	return map[string]interface{}{"success": true, "status": "Filter Refreshed"}, nil
}

func (s *Storage) removeFilterValues(data map[string]interface{}) bool {
	log.Printf("Cleaning old filter values for domain: %d", data["domainId"])

	deleteSQL := `
		DELETE FROM sch_filter_values
		WHERE sch_schema_id = $1
		AND metadata_id IN (
			SELECT sm.sch_metadata_id 
			FROM sch_metadata sm
			WHERE sm.sch_schema_id = $1
			AND (sm.enable_for_filter = 0 OR sm.is_active = 0 OR sm.column_is_active = 0)
		)
	`

	if _, err := s.Db.Exec(deleteSQL, data["domainId"]); err != nil {
		log.Printf("Error cleaning old filter values: %v", err)
		return true
	}
	return true
}

func (s *Storage) fetchMetadataIDs(domainID interface{}) ([]int64, error) {
	query := `
		SELECT sm.sch_metadata_id
		FROM sch_metadata sm
		WHERE sm.sch_schema_id = ? 
		AND sm.enable_for_filter = 1 
		AND sm.is_active 
		AND sm.column_is_active 
		AND sm.column_type IN ('Dimension', 'Location', 'vDimension')
		AND sm.sch_metadata_id NOT IN 
		(SELECT DISTINCT svf.metadata_id FROM sch_filter_values svf WHERE svf.sch_schema_id = ?)
	`

	rows, err := s.QueryWithTiming(query, domainID, domainID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch metadata IDs: %w", err)
	}

	metadataIds, err := GetFromRows(rows)
	if err != nil {
		fmt.Printf("unable to fetch metadata ids error %v", err.Error())
	}

	var _metadataIds []int64
	for _, meta := range metadataIds {
		id, ok := meta["sch_metadata_id"].(int64)
		if !ok {
			fmt.Printf("metadata id not found")
		}
		_metadataIds = append(_metadataIds, id)
	}

	return _metadataIds, nil
}

func fetchFilterValues(s *Storage, domainID int, metadataIDs []int64, filterLimit int64) (any, error) {
	var wg sync.WaitGroup

	type Filter struct {
		FilterName  string `json:"filter_name"`
		MetadataID  int64  `json:"metadata_id"`
		SchSchemaID int    `json:"sch_schema_id"`
	}

	filterCh := make(chan Filter, len(metadataIDs))

	for _, metaId := range metadataIDs {
		wg.Add(1)

		go func(metaId int64) {

			defer wg.Done()
			filters, err := getDistinctValues(s, domainID, metaId, filterLimit)

			if err != nil {
				fmt.Printf("unable to fetch filter from analytical service %v\n", err.Error())
			}

			for _, filter := range filters {
				filterCh <- Filter{FilterName: filter, MetadataID: metaId, SchSchemaID: domainID}

			}

		}(metaId)
	}

	go func() {
		wg.Wait()
		close(filterCh)
	}()

	var filters []Filter
	for filter := range filterCh {
		filters = append(filters, filter)
	}

	return filters, nil
}

func getDistinctValues(s *Storage, domainID int, metadataID, limit int64) ([]string, error) {
	var filters []string
	payload := map[string]interface{}{
		"schemaId":   domainID,
		"moduleName": "askme-queue",
		"attributes": []map[string]interface{}{
			{
				"formula":    "count(1)",
				"columnType": "ac",
				"langId":     3,
				"aliasName":  "TotalRows",
				"sort":       "DESC",
				"table":      []any{},
			},
			{
				"attribute": metadataID,
			},
		},
		"limit": limit,
	}

	url := s.Parameters["analyticalLayer"].(string) + "/get-query"

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	_filters, err := s.MakeRequest("POST", url, payload, headers)
	if err != nil {
		log.Printf("error while fetching filter for metadataId %d: err %v\n", metadataID, err.Error())
	}

	if filters, ok := (*_filters).([]string); ok {
		fmt.Println("Successfully converted to []string:", filters)
	} else {
		fmt.Println("Failed to convert to []string")
	}

	return filters, nil

}

func (s *Storage) insertFilters(domainID, userID interface{}, filters any) error {
	// insertSQL := "INSERT INTO sch_filter_values (filter_name, metadata_id, sch_schema_id, modified_by) VALUES "
	// var builder strings.Builder

	// for _, filter := range filters {
	// 	builder.WriteString(fmt.Sprintf(`('%s', %d, %d, %d),`, filter.FilterName, filter.MetadataID, domainID, userID))
	// }

	// query := builder.String()
	// if len(query) > 0 {
	// 	query = query[:len(query)-1]
	// 	_, err := s.Db.Exec(query)
	// 	if err != nil {
	// 		log.Printf("Error inserting filters: %v", err)
	// 		return fmt.Errorf("error inserting filters: %w", err)
	// 	}
	// }

	return nil
}

func (s *Storage) getFilterLimit(data map[string]interface{}) int64 {
	filterLimit := int64(500)
	query := `
		SELECT data_limit
		FROM sch_analytical_data_limit
		WHERE is_active = 1
			AND module_name = 'askme-filter-limit'
			AND data_limit IS NOT NULL
			AND (
				org_id = -1
				OR (
					org_id = ?
					AND (user_id = ? OR user_id IS NULL)
					AND (sch_schema_id = ? OR sch_schema_id IS NULL)
				)
			)
		ORDER BY org_id DESC, sch_schema_id DESC, user_id DESC
		LIMIT 1;
	`

	var result int64
	if err := s.Db.QueryRow(query, data["orgID"], data["userID"], data["domainID"]).Scan(&result); err != nil {
		log.Printf("Error fetching filter limit: %v", err)
		return filterLimit
	}

	return result
}
