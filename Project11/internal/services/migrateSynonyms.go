package services

import (
	"log"
	"strconv"
)

func (s *Storage) MigrateSynonyms(data map[string]interface{}) (any, error) {
	log.Println("Migrating Synonyms")

	// rows, err := s.Db.Query(`
	// SELECT type_id, synonym as synonym,  asn.is_global, asn.user_id, word_type as word_type, column_alias as word,
	//     CASE
	//         WHEN ('column_type' = 'Measure') THEN 'aM'
	//         WHEN ('column_type' = 'vMeasure') THEN 'avM'
	//         WHEN ('column_type' = 'Dimension' and 'column_is_date' = 1) THEN 'aDD'
	//         WHEN ('column_type' = 'Dimension') THEN 'aD'
	//         WHEN ('column_type' = 'vDimension') THEN 'avD'
	//         WHEN ('column_type' = 'Calculated KPI') THEN 'aC'
	//         WHEN ('column_type' = 'vCalculated KPI') THEN 'avC'
	//         WHEN ('column_type' = 'Location') THEN 'aLoc'
	//     END as category
	//     FROM ask_synonym asn join sch_metadata sm
	//     ON asn.type_id = sm.sch_metadata_id
	//     WHERE sm.column_is_active = 1 and  sm.is_active = 1 and sm.sch_schema_id = ? and asn.is_active= 1
	//     UNION ALL
	//     SELECT type_id,  synonym as synonym, asn.is_global, asn.user_id, word_type as word_type, keyword_name as word ,
	//     CASE
	//         WHEN ( 'keyword_category' = 'Condition') THEN 'aCondition'
	//         WHEN ( 'keyword_category' = 'Aggregated Function') THEN 'aAM'
	//         WHEN ( 'keyword_category' = 'Arrange') THEN 'aTB'
	//         WHEN ( 'keyword_category' = 'Date Group') THEN 'aDate'
	//         WHEN ( 'keyword_category' = 'Function') THEN 'aAD'
	//         WHEN ( 'keyword_category' = 'Ordering') THEN 'aGF'
	//         WHEN ( 'keyword_category' = 'Date Column Filter' ) THEN 'aDF'
	//         WHEN ( 'keyword_category' = 'Date Filter' ) THEN 'aTDF'
	//         WHEN ( 'keyword_category' = 'Rank' ) THEN 'aR'
	//         WHEN ( 'keyword_category' = 'Change Analysis' ) THEN 'aG'
	//         WHEN ( 'keyword_category' = 'Null' ) THEN 'aNULL'
	//         WHEN ( 'keyword_category' = 'Advanced Analysis' ) THEN 'aAA'
	//         WHEN ( 'keyword_category' = 'Order by' ) THEN 'aOB'
	//         WHEN ( 'keyword_category' = 'Order' ) THEN 'aO'
	//     END as category
	//     from ask_synonym asn join ask_keyword ak
	//     on asn.type_id = ak.ask_keyword_id
	//     WHERE ak.is_active = 1 and asn.is_active= 1 and asn.sch_schema_id = ?
	//     UNION ALL
	//     SELECT sch_metadata_id as type_id, synonym as synonym,  asn.is_global, asn.user_id, word_type as word_type,
	//     filter_name as word , sm.column_alias as category
	//     from ask_synonym asn join sch_filter_values af
	//     on asn.type_id = af.sch_filter_value_id
	//     join sch_metadata sm
	//     on af.metadata_id = sm.sch_metadata_id
	//     WHERE sm.column_is_active = 1 and sm.is_active = 1 and sm.enable_for_filter=1 and
	//     sm.sch_schema_id = ? and asn.is_active= 1`, data["domainId"], data["domainId"], data["domainId"])

	rows, err := s.QueryWithTiming(`
				with metadata as (select sch_metadata_id, column_alias, column_type, column_is_date from sch_metadata sm where sm.sch_schema_id = ?),
			askfilter as (select sfv.* from sch_filter_values sfv inner join metadata md on sfv.metadata_id = md.sch_metadata_id),
			synonyms as (select type_id, synonym, is_global, user_id, word_type from ask_synonym where sch_schema_id = ? and is_active = 1),
			askkeyword as (select * from ask_keyword ak where is_active = 1)
			select
				syn.type_id,
				syn.synonym as synonym,
				syn.is_global,
				syn.user_id,
				syn.word_type,
				case
					when syn.word_type = 'metadata' then sm.column_alias
					when syn.word_type = 'keyword' then ak.keyword_name
					else af.filter_name
				end as word,
				case
					when syn.word_type = 'metadata' then
						case
							when (sm.column_type = 'Measure') then 'aM'
							when (sm.column_type = 'vMeasure') then 'avM'
							when (sm.column_type = 'Dimension' and sm.column_is_date = 1) then 'aDD'
							when (sm.column_type = 'Dimension') then 'aD'
							when (sm.column_type = 'vDimension') then 'avD'
							when (sm.column_type = 'Calculated KPI') then 'aC'
							when (sm.column_type = 'vCalculated KPI') then 'avC'
							when (sm.column_type = 'Location') then 'aLoc'
						end
					when syn.word_type = 'keyword' then
						case
							when ( ak.keyword_category = 'Condition') then 'aCondition'
							when ( ak.keyword_category = 'Aggregated Function') then 'aAM'
							when ( ak.keyword_category = 'Arrange') then 'aTB'
							when ( ak.keyword_category = 'Date Group') then 'aDate'
							when ( ak.keyword_category = 'Function') then 'aAD'
							when ( ak.keyword_category = 'Ordering') then 'aGF'
							when ( ak.keyword_category = 'Date Column Filter' ) then 'aDF'
							when ( ak.keyword_category = 'Date Filter' ) then 'aTDF'
							when ( ak.keyword_category = 'Rank' ) then 'aR'
							when ( ak.keyword_category = 'Change Analysis' ) then 'aG'
							when ( ak.keyword_category = 'Null' ) then 'aNULL'
							when ( ak.keyword_category = 'Advanced Analysis' ) then 'aAA'
							when ( ak.keyword_category = 'Order by' ) then 'aOB'
							when ( ak.keyword_category = 'Order' ) then 'aO'
						end
					else
						sm.column_alias
				end as category
			from synonyms syn
			left join sch_metadata sm on syn.type_id = sm.sch_metadata_id
			left join askkeyword ak on syn.type_id = ak.ask_keyword_id
			left join askfilter af on syn.type_id = af.sch_filter_value_id
			having word is not null and category is not null`, data["domainId"], data["domainId"])
	if err != nil {
		return nil, err
	}

	synonym, err := GetFromRows(rows)
	if err != nil {
		return nil, err
	}
	var _syn = make(map[string]interface{})

	synonymDict := map[string]map[string]interface{}{
		"global": {},
		"user":   {},
	}

	for _, syn := range synonym {
		synonymStr := typeCases(syn["synonym"])

		key := cleanQuery(synonymStr)

		if syn["is_global"].(uint8) == 1 {
			synonymDict["global"][key] = syn

			if syn["word_type"].(string) == "metadata" {
				if _, exists := synonymDict["user"]["global_syn"]; !exists {
					synonymDict["user"]["global_syn"] = map[string]interface{}{}
				}
				synonymDict["user"]["global_syn"].(map[string]interface{})[key] = syn
			}
		} else {
			userID := typeCases(syn["user_id"])
			delete(syn, "user_id")

			if _, exists := synonymDict["user"][userID]; !exists {
				synonymDict["user"][userID] = map[string]interface{}{}
			}
			synonymDict["user"][userID].(map[string]interface{})[key] = syn
		}
	}
	// domainIDStr := fmt.Sprintf("%0.f", data["domainId"])
	log.Printf("_syn %+v", _syn)

	// go s.hPublish(domainIDStr, "synonyms", _syn)
	return synonymDict, nil

}

func typeCases(syn any) string {
	var synonymStr string

	switch v := syn.(type) {
	case string:
		synonymStr = v
	case []byte:
		synonymStr = string(v)
	case uint8:
		synonymStr = string(v)
	case int:
		synonymStr = strconv.Itoa(v)
	case int64:
		synonymStr = strconv.FormatInt(v, 10)
	default:
		log.Printf("Error: synonym is of unexpected type %T", v)
	}
	return synonymStr
}
