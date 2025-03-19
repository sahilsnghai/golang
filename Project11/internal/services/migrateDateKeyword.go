package services

import (
	"fmt"
	"log"
)

func (s *Storage) MigrateDateKeyword(data map[string]interface{}) (any, error) {

	log.Println("Migrating Date Keyword")

	rows, err := s.QueryWithTiming(`select
        ask_keyword_id as type_id,
        'keyword' as word_type,
        upper( keyword_name ) as word,
        CASE
            WHEN ( keyword_category = 'Date Column Filter' ) THEN 'aDF'
        END AS category
    from
        ask_keyword ak
    where
        upper( keyword_name ) = 'DAY'`)

	if err != nil {
		return nil, fmt.Errorf("uable to fetch aDF %s", err.Error())
	}

	aDFs, err := GetFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("uable to retrive rows for aDF %s", err.Error())
	}
	aDF := aDFs[0]
	go s.hPublish(fmt.Sprintf("%0.f", data["domainId"]), "aDF", aDF)

	return aDF, nil
}
