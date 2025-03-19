package services

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

func (s *Storage) refreshCahe(wg *sync.WaitGroup) (any, error) {

	fmt.Print("refreshCache")
	lst := map[string]func() (interface{}, error){
		"askme_keywords":         s.GetKeywords,
		"askme_keyword_synonyms": s.GetKeywordSynonyms,
		"askme_keyword_names":    s.GetKeywordName,
	}

	var caches = make(map[string]any)

	for key, fn := range lst {
		wg.Add(1)
		defer wg.Done()
		cache, err := fn()
		if err != nil {
			log.Printf("error while fetching keywords for %s: %s", key, err.Error())
			return nil, err
		}
		caches[key] = cache

		go s.Publish("", key, cache)
	}

	return caches, nil

}

func (s *Storage) GetKeywords() (interface{}, error) {
	log.Println("GetKeywords")

	rows, err := s.QueryWithTiming(`
				 SELECT
					ask_keyword_id,
					keyword_name,
					keyword_formula,
					ask_keyword_alias as keyword_category
				FROM
					ask_keyword
				where
					is_active = 1;`,
	)
	if err != nil {
		return nil, err
	}

	keywords, err := GetFromRows(rows)
	if err != nil {
		return nil, err
	}

	kWords := make(map[string]map[string]interface{})
	for _, item := range keywords {
		word, ok := item["ask_keyword_id"].(int64)
		if ok {
			delete(item, "ask_keyword_id")
			upperWord := strings.ToUpper(strconv.Itoa(int(word)))
			kWords[upperWord] = item
		}
	}

	return kWords, nil

}

func (s *Storage) GetKeywordSynonyms() (interface{}, error) {
	log.Println("GetKeywordSynonyms")

	// var keywordSynonyms []KeywordSynonym

	// err := s.Db2.Table("ask_synonym asn").
	// 	Select("asn.type_id, asn.synonym, asn.word_type, ak.keyword_name as word, ak.ask_keyword_alias as category").
	// 	Joins("JOIN ask_keyword ak ON asn.type_id = ak.ask_keyword_id").
	// 	Where("ak.is_active = ? AND asn.is_active = ? AND asn.sch_schema_id = ? AND asn.word_type = ?",
	// 		1, 1, 1, "keyword").
	// 	Scan(&keywordSynonyms).Error

	rows, err := s.QueryWithTiming(`SELECT
					type_id,
					synonym AS synonym,
					word_type AS word_type,
					keyword_name AS word ,
					ak.ask_keyword_alias AS category
				FROM
						ask_synonym asn
				JOIN ask_keyword ak ON
						asn.type_id = ak.ask_keyword_id
				WHERE
						ak.is_active = ?
						AND asn.is_active = ?
						AND asn.sch_schema_id = ?
						AND asn.word_type = ?`,
		1, 1, 1, "keyword")

	defer rows.Close()
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}

	keywordSynonyms, err := GetFromRows(rows)
	if err != nil {
		return nil, err
	}

	synonymMap := make(map[string]interface{})
	for _, value := range keywordSynonyms {
		var synonymStr string

		switch v := value["synonym"].(type) {
		case string:
			synonymStr = v
		case []byte:
			synonymStr = string(v)
		case uint8:
			synonymStr = string(v)
		default:
			log.Printf("Error: synonym is of unexpected type %T", v)
			continue
		}
		synonym := strings.ToUpper(synonymStr)
		synonymMap[synonym] = value
	}

	return map[string]map[string]interface{}{
		"global": synonymMap,
		"user":   make(map[string]interface{}),
	}, nil
}

func (s *Storage) GetKeywordName() (interface{}, error) {
	log.Println("GetKeywordName")

	rows, err := s.QueryWithTiming(`select
					ask_keyword_id as type_id,
					'keyword' as word_type,
					keyword_name as word,
					ask_keyword_alias AS category
				from
					ask_keyword ak
				where
					ak.is_active = 1`,
	)
	if err != nil {
		return nil, err
	}

	keyWordName, err := GetFromRows(rows)
	if err != nil {
		return nil, err
	}

	keyWords := make(map[string][]map[string]interface{})
	for _, value := range keyWordName {
		word := strings.ToUpper(value["word"].(string))

		if _, exists := keyWords[word]; !exists {
			keyWords[word] = []map[string]interface{}{}
		}
		keyWords[word] = append(keyWords[word], value)
	}

	return keyWords, nil
}
