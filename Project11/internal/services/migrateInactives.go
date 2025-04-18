package services

import (
	"fmt"
	"log"
)

func (s *Storage) MigrateInactive(data map[string]interface{}) (any, error) {
	log.Println("Migrate Inactive")

	domainID, ok := data["domainId"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing domainId")
	}

	// var inactives []InActiveColumns
	// err := s.Db2.Select("sm.column_alias").
	// 	Table("sch_metadata sm").
	// 	Where("sm.sch_schema_id = ? AND NOT sm.column_is_active", domainID).
	// 	Scan(&inactives).Error

	rows, err := s.QueryWithTiming(`select sm.column_alias from sch_metadata sm where sm.sch_schema_id = ? AND NOT sm.column_is_active`, domainID)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch inactive columns: %w", err)
	}

	inactives, err := GetFromRows(rows)

	if err != nil {
		return nil, fmt.Errorf("failed to retrive rows for Inactive: %w", err)
	}

	inactiveSlice := make([]string, 0, len(inactives))
	for _, item := range inactives {
		inactiveSlice = append(inactiveSlice, item["column_alias"].(string))
	}

	s.hPublish(fmt.Sprintf("%0.f", domainID), "inactive", inactiveSlice)

	return inactiveSlice, nil
}
