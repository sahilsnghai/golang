package services

type KeywordSynonym struct {
	TypeID          int    `gorm:"column:type_id"`
	Synonym         string `gorm:"column:synonym"`
	WordType        string `gorm:"column:word_type"`
	KeywordName     string `gorm:"column:keyword_name"`
	AskKeywordAlias string `gorm:"column:ask_keyword_alias"`
}

type Keyword struct {
	AskKeywordId    int    `gorm:"column:ask_keyword_id"`
	KeywordName     string `gorm:"column:keyword_name"`
	KeywordFormula  string `gorm:"column:keyword_formula"`
	AskKeywordAlias string `gorm:"column:ask_keyword_alias"`
}

type InActiveColumns struct {
	ColumnAlias string `grom:"column:column_alias"`
}

type BasicJoins struct {
	SchJoinTableId  int    `gorm:"column:sch_join_table_id"`
	JoinType        string `gorm:"column:join_type"`
	LeftTableName   string `gorm:"column:left_table_name"`
	RightTableName  string `gorm:"column:right_table_name"`
	LeftColumnName  string `gorm:"column:left_column_name"`
	RightColumnName string `gorm:"column:right_column_name"`
	JoinCondition   string `gorm:"column:join_condition"`
	ConditionMerger string `gorm:"column:condition_merger"`
}
