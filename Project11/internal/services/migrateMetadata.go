package services

import (
	"fmt"
	"log"
	"strconv"
)

func (s *Storage) MigrateMetadata(data map[string]interface{}) (any, error) {
	log.Println("Migrating Metadata")

	rows, err := s.QueryWithTiming(`
    SELECT
            sm.sch_metadata_id AS metadata_id,
            sm.column_name AS metadata_column_name,
            sm.column_alias AS metadata_column_alias,
            sm.formula AS metadata_formula,
            UPPER(sm.column_table_alias) AS metadata_column_table_alias,
            sm.column_table_name AS metadata_column_table_name,
            CASE
                WHEN sm.column_type = 'Measure' THEN 'aM'
                WHEN sm.column_type = 'vMeasure' THEN 'avM'
                WHEN (sm.column_type = 'Dimension' AND sm.column_is_date = 1) THEN 'aDD'
                WHEN (sm.column_type = 'vDimension' AND sm.column_is_date = 1) THEN 'avD'
                WHEN sm.column_type = 'Dimension' THEN 'aD'
                WHEN sm.column_type = 'vDimension' THEN 'avD'
                WHEN sm.column_type = 'Calculated KPI' THEN 'aC'
                WHEN sm.column_type = 'vCalculated KPI' THEN 'avC'
                WHEN sm.column_type = 'Latitude' THEN 'aLat'
                WHEN sm.column_type = 'Longitude' THEN 'aLng'
                WHEN sm.column_type = 'Location' THEN 'aLoc'
                WHEN sm.column_type = 'Variable' THEN 'aVar'
            END AS metadata_column_type,
            lrf_query_lang_id AS metadata_column_lang_id,
            sm.column_is_date AS metadata_column_is_date,
            cdf.date_format AS metadata_column_date_format,
            ak.keyword_formula as metadata_default_agg_type_formula,
            CASE
                WHEN ak.keyword_name = 'NA' THEN NULL
                ELSE ak.keyword_name
            END AS metadata_default_agg_type,
            sm.column_datatype_name AS metadata_column_datatype_name,
            sm.column_custom_datatype AS metadata_custom_datatype,
            CASE
                WHEN sm.column_default_order = 'NA' THEN NULL
                ELSE sm.column_default_order
            END AS metadata_column_ordering,
            sm.column_is_encrypted AS is_encrypted,
            sm.is_auto_generated AS metadata_auto_generated,
            sm.unit AS metadata_unit,
            sm.unit_placement AS metadata_unit_placement,
            sm.is_what_if AS metadata_is_what_if,
            sm.column_is_pre_calculated AS metadata_pre_calculated,
            sm.what_if_max AS metadata_what_if_max,
            sm.what_if_min AS metadata_what_if_min,
            sm.what_if_step AS metadata_what_if_step,

            CASE
            WHEN
                (POSITION('analytical' IN sm.formula_ast) > 0 OR
                 POSITION('aggregate_analytical' IN sm.formula_ast) > 0 OR
                 POSITION('aggregate_analytical_min' IN sm.formula_ast) > 0 OR
                 POSITION('aggregate_analytical_max' IN sm.formula_ast) > 0 OR
                 POSITION('aggregate_analytical_stddev' IN sm.formula_ast) > 0 OR
                 POSITION('aggregate_analytical_var' IN sm.formula_ast) > 0 OR
                 POSITION('aggregate_analytical_count' IN sm.formula_ast) > 0 OR
                 POSITION('groupby_aggregate' IN sm.formula_ast) > 0 OR
                 POSITION('groupby_aggregate_include_attr' IN sm.formula_ast) > 0 OR
                 POSITION('groupby_aggregate_exinclude_filter' IN sm.formula_ast) > 0 OR
                 POSITION('groupby_aggregate_include_attr_exinclude_filter' IN sm.formula_ast) > 0 OR
                 POSITION('analytical' IN sm.formula) > 0 OR
                 POSITION('aggregate_analytical' IN sm.formula) > 0 OR
                 POSITION('aggregate_analytical_min' IN sm.formula) > 0 OR
                 POSITION('aggregate_analytical_max' IN sm.formula) > 0 OR
                 POSITION('aggregate_analytical_stddev' IN sm.formula) > 0 OR
                 POSITION('aggregate_analytical_var' IN sm.formula) > 0 OR
                 POSITION('aggregate_analytical_count' IN sm.formula) > 0 OR
                 POSITION('groupby_aggregate' IN sm.formula) > 0 OR
                 POSITION('groupby_aggregate_include_attr' IN sm.formula) > 0 OR
                 POSITION('groupby_aggregate_exinclude_filter' IN sm.formula) > 0 OR
                 POSITION('groupby_aggregate_include_attr_exinclude_filter' IN sm.formula) > 0)
            THEN FALSE
            ELSE TRUE
            END AS is_formula_valid,

            sm.color AS metadata_color,
            sm.thousand_separator AS metadata_thousand_separator,
            sm.formatting_type AS metadata_formatting_type,
            sm.precision AS metadata_precision,
            sm.abbreviation AS metadata_abbreviation,
            sm.formatter_text_position AS metadata_formatter_text_position,
            sm.description AS metadata_description

        FROM  sch_metadata sm
        LEFT JOIN sch_column_date_format cdf ON
            sm.sch_schema_id = cdf.sch_schema_id
            AND
            sm.sch_metadata_id = cdf.metadata_id
        LEFT JOIN ask_keyword ak ON sm.default_agg_type_id = ak.ask_keyword_id
        WHERE sm.sch_schema_id = ?
            AND sm.column_is_active = TRUE
            AND sm.is_active = 1
            AND sm.is_global = 1;
    `, data["domainId"])

	if err != nil {
		return nil, err
	}

	metadata, err := GetFromRows(rows)
	if err != nil {
		return nil, err
	}

	var meta = make(map[string]interface{})
	for _, _meta := range metadata {
		_id := strconv.FormatInt(_meta["metadata_id"].(int64), 10)
		meta[_id] = _meta
	}
	domainIDStr := fmt.Sprintf("%0.f", data["domainId"])

	go s.hPublish(domainIDStr, "metadata", meta)

	return meta, nil

}
