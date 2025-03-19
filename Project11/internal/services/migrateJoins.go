package services

import (
	"fmt"
	"log"
)

func (s *Storage) MigrateJoins(data map[string]interface{}) (any, error) {
	log.Println("Migrating Basic Joins")

	domainID, ok := data["domainId"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing domainId")
	}
	rows, err := s.QueryWithTiming(`Select sjt.sch_join_table_id, join_type, left_table_name, right_table_name,
        left_column_name, right_column_name, join_condition, condition_merger
    from
        sch_join_table sjt
    inner join sch_join_condition sjc on
        sjt.sch_join_table_id = sjc.sch_join_table_id
    where
        sch_schema_id = ?
        and sjt.is_active = ?
        and sjc.is_active = ?`, domainID, 1, 1)

	if err != nil {
		return nil, fmt.Errorf("error cannot fetch basic joins from the table: %w", err)
	}

	basicJoin, err := GetFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("not able to convert into slices of maps %w", err)
	}
	var JoinsResults = make(map[int64]map[string]interface{})

	for _, basicjn := range basicJoin {
		SchJoinId := basicjn["sch_join_table_id"].(int64)
		var conditions = map[string]string{
			"left_column_name":  string(basicjn["left_column_name"].(string)),
			"right_column_name": string(basicjn["right_column_name"].(string)),
			"condition":         string(basicjn["join_condition"].(byte)),
			"condition_merger":  string(basicjn["condition_merger"].(string)),
		}

		if _, exists := JoinsResults[SchJoinId]; exists {
			JoinsResults[SchJoinId]["conditions"] = append(JoinsResults[SchJoinId]["conditions"].([]map[string]string), conditions)
		} else {
			JoinsResults[SchJoinId] = map[string]interface{}{
				"join_type":        string(basicjn["join_type"].(string)),
				"left_table_name":  string(basicjn["left_table_name"].(string)),
				"right_table_name": string(basicjn["right_table_name"].(string)),
				"conditions":       []map[string]string{conditions},
			}
		}
	}

	JoinsValues := func(m map[int64]map[string]interface{}) []interface{} {
		k := make([]interface{}, 0, len(m))
		for _, value := range m {
			k = append(k, value)
		}
		return k
	}(JoinsResults)

	go s.hPublish(fmt.Sprintf("%0.f", data["domainId"]), "basic_joins", JoinsValues)

	return JoinsValues, nil

}
