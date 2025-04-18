package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type MappingEntry struct {
	Process func(data map[string]interface{}) (any, error)
	Params  []string
}

func (s Storage) Migration(data map[string]interface{}) (map[string]interface{}, error) {
	log.Println("migration")
	var wg sync.WaitGroup
	defer wg.Wait()
	s.Ctx = context.Background()
	var _metadata map[string]interface{}

	ex1, err1 := s.Client.Exists(s.Ctx, "askme_keywords").Result()
	ex2, err2 := s.Client.Exists(s.Ctx, "askme_keyword_synonyms").Result()
	if (ex1 == 0 || err1 != nil) || (ex2 == 0 || err2 != nil) {
		log.Printf("Either 'askme_keywords' or 'askme_keyword_synonyms' does not exist, or there was an error. %v , %v", ex1, ex2)
		go s.refreshCahe(&wg)
	}

	// filter-refresh true
	if filtersRef, ok := data["filter-refresh"].(bool); ok && filtersRef {
		log.Printf("Calling filter refresh")
		s.FilterRefresh(data)
	}
	log.Printf("Request data from for migration: %+v\n", data)

	// migrate-metadata true
	if migrateMetadata, ok := data["migrate-metadata"].(bool); ok && migrateMetadata {
		log.Printf("Calling MIgrate Metadata")

		var err error
		_metadata, err = s.Migrate(data, &wg)
		if err != nil {
			log.Printf("error while migration: %v ", err.Error())
			return nil, err
		}
	}

	return _metadata, nil
}

func (s *Storage) Migrate(data map[string]interface{}, wg *sync.WaitGroup) (map[string]interface{}, error) {

	log.Printf("starting migration")
	domainIDStr := fmt.Sprintf("%0.f", data["domainId"])

	mappings := map[string]MappingEntry{
		"basic_joins":       {s.MigrateJoins, []string{"domain_id", "hash"}},       // Done implementation
		"words_suggestions": {s.MigrateWords, []string{"domain_id", "hash"}},       // Done implementation
		"credentials":       {s.MigrateCreds, []string{"domain_id", "hash"}},       // Done implementation
		"inactive":          {s.MigrateInactive, []string{"domain_id", "hash"}},    // Done implementation
		"metadata":          {s.MigrateMetadata, []string{"domain_id", "hash"}},    // Done implementation
		"synonyms":          {s.MigrateSynonyms, []string{"domain_id", "hash"}},    // Done implementation
		"hierarchy":         {s.MigrateHierarchy, []string{"domain_id", "hash"}},   // Done implementation
		"shortcuts":         {s.MigrateShortCuts, []string{"domain_id", "hash"}},   // Done implementation
		"aDD":               {s.MigrateDateColumn, []string{"domain_id", "hash"}},  // Done implementation
		"aDF":               {s.MigrateDateKeyword, []string{"domain_id", "hash"}}, // Done implementation
		// "sample":            {processSample, []string{"domain_id", "main_dict", "redis"}},
	}

	key, ok := data["key"].(string)
	log.Printf("key: %s, domainIDStr: %s\n", key, domainIDStr)

	var migrateMappings map[string]MappingEntry
	if !ok {
		migrateMappings = mappings

	} else if key == "metadata" || key == "derived_kpi" {
		migrateMappings = map[string]MappingEntry{
			"metadata":          {s.MigrateMetadata, []string{"domain_id", "redis"}},
			"aDD":               {s.MigrateDateColumn, []string{"domain_id", "user_id", "redis"}},
			"inactive":          {s.MigrateInactive, []string{"domain_id", "redis"}},
			"words_suggestions": {s.MigrateWords, []string{"domain_id", "redis"}},
			// "sample",
		}

	} else {
		migrateMappings = map[string]MappingEntry{
			key: mappings[key],
		}
	}
	log.Printf("migrateMapping %+v \n", migrateMappings)
	result := make(map[string]interface{})

	var mu sync.Mutex

	startTime := time.Now()
	for key, fCall := range mappings {
		wg.Add(1)
		go func(k string, f MappingEntry) {
			defer wg.Done()

			log.Printf("Starting migration: %s\n", k)
			if f.Process != nil {
				_data, err := f.Process(data)
				if err != nil {
					log.Printf("Error while migrating %s: %s\n", k, err.Error())
					return
				}
				mu.Lock()
				result[k] = _data
				mu.Unlock()

				log.Printf("Completed migration: %s\n", k)
			}
		}(key, fCall)
	}
	elapsedTime := time.Since(startTime)
	log.Printf("\tComplete metadata fetch in %v\n", elapsedTime)
	return result, nil

}
