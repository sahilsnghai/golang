package services

import (
	"fmt"
	"log"
)

func (s *Storage) MigrateHierarchy(data map[string]interface{}) (any, error) {
	log.Println("Migrating Hierarchy")

	rows, err := s.QueryWithTiming(`SELECT ah.user_id,
				ah.is_global,
				UPPER(sm1.column_alias) AS from_column_name,
				UPPER(sm2.column_alias) AS to_column_name
			FROM
				ask_hierarchy_detail ahdt
			INNER JOIN ask_hierarchy ah ON
				ahdt.ask_hierarchy_id = ah.ask_hierarchy_id
			INNER JOIN sch_metadata sm1 ON
				ahdt.from_metadata = sm1.sch_metadata_id
			INNER JOIN sch_metadata sm2 ON
				ahdt.to_metadata = sm2.sch_metadata_id
			WHERE
				ah.sch_schema_id = ?
				AND ah.is_active
				AND ahdt.is_active`, data["domainId"])

	if err != nil {
		return nil, fmt.Errorf("uable to fetch Hierarchy %s", err.Error())
	}

	hierarchys, err := GetFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("unable to retrive rows for Hierarchy %s", err.Error())
	}
	var hierarchy = map[string]interface{}{
		"default": [][]string{},
		"global":  [][]string{},
		"user":    make(map[int64][][]string),
	}

	for _, _hrchy := range hierarchys {
		userId := _hrchy["user_id"].(int64)

		if isGlobal, ok := _hrchy["is_global"].(uint8); ok && isGlobal == 1 {
			hierarchy["global"] = append(hierarchy["global"].([][]string), []string{
				_hrchy["from_column_name"].(string),
				_hrchy["to_column_name"].(string),
			})
		}
		userMap := hierarchy["user"].(map[int64][][]string)
		userMap[userId] = append(userMap[userId], []string{
			_hrchy["from_column_name"].(string),
			_hrchy["to_column_name"].(string),
		})
	}

	go s.hPublish(fmt.Sprintf("%0.f", data["domainId"]), "hierarchy", hierarchy)

	return hierarchy, nil
}
