package services

import (
	"fmt"
	"log"
	"strings"
)

func (s *Storage) MigrateShortCuts(data map[string]interface{}) (any, error) {
	log.Println("Migrating Shortcuts")

	rows, err := s.QueryWithTiming(`
    SELECT
		ask_shortcut_id,
		name,
		shortcut_text,
		attribute_info,
		is_global,
		modified_by
	FROM
		ask_shortcuts
	WHERE
        sch_schema_id = ? and is_active`,
		data["domainId"])

	if err != nil {
		return nil, err
	}

	shortCuts, err := GetFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("unable to retrive rows for shortcuts %s", err.Error())
	}

	localSC := make(map[string]interface{})
	globalSC := make(map[string]interface{})
	completeSC := make(map[string]interface{}, 0)

	for _, shorts := range shortCuts {

		name := strings.ToUpper(shorts["name"].(string))
		modifiedBy := shorts["modified_by"].(int64)

		if isGlobal, ok := shorts["is_global"].(uint8); ok && isGlobal == 1 {
			globalSC[name] = shorts
		} else if modifiedBy == data["userId"] {
			localSC[name] = shorts
		}
		completeSC[name] = shorts

	}

	go s.hPublish(fmt.Sprintf("%0.f", data["domainId"]), "shortcuts", globalSC)
	go s.Publish(fmt.Sprintf("%0.f", data["domainId"]),
		fmt.Sprintf("amsc_%0.f_%d", data["domainId"],
			data["userId"]), globalSC)

	return completeSC, nil
}
