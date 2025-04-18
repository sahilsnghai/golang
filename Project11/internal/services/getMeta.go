package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *Storage) GetMetadata(getMetadata map[string]interface{}) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	var wg sync.WaitGroup
	var mu sync.Mutex

	keys, ok := getMetadata["keys"].([]interface{})
	if !ok || len(keys) == 0 {
		keys = []interface{}{"metadata", "words", "hierarchy", "inactive", "keyword", "aDD", "aDF", "synonyms", "shortcuts"}
	}

	userID, ok := getMetadata["userId"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid type for req, expected userId")
	}

	domainIDStr := fmt.Sprintf("%0.f", getMetadata["domainId"])
	start1 := time.Now()
	ctx := context.Background()

	for _, v := range keys {
		key, ok := v.(string)
		if !ok {
			continue
		}
		wg.Add(1)
		log.Printf("Processing key: %s, domainId: %s, userId: %0.f\n", key, domainIDStr, userID)

		go func(key string) {
			defer wg.Done()
			value, err := fetchMetadataValue(key, domainIDStr, userID, s.Client, ctx, s)
			if err != nil {
				log.Printf("Error processing key %s: %v", key, err)
			}

			mu.Lock()
			results[key] = value
			mu.Unlock()
		}(key)
	}

	wg.Wait()
	log.Printf("Total time for getMetadata: %v ms", time.Since(start1).Milliseconds())
	return results, nil
}

func fetchMetadataValue(key, domainIDStr string, userID float64, client *redis.Client, ctx context.Context, s *Storage) (interface{}, error) {
	var value interface{}
	var err error
	switch key {
	case "words", "suggestion":
		value, err = getWords(domainIDStr, client, ctx)
		if err != nil {
			return nil, fmt.Errorf("error getting words: %v", err)
		}
	case "keyword":
		value = client.Get(ctx, "askme_keywords").Val()
	case "synonyms":
		value, err = getSynonyms(domainIDStr, client, ctx, s)
		if err != nil {
			return nil, fmt.Errorf("error getting synonyms: %v", err)
		}
	case "shortcuts":
		value, err = getShortcuts(domainIDStr, userID, client, ctx)
		if err != nil {
			return nil, fmt.Errorf("error getting shortcuts: %v", err)
		}
	default:
		value = fetchRedisHashValue(key, domainIDStr, client, ctx)
	}

	if rawValue, ok := value.(string); ok && rawValue != "" {
		var tempValue interface{}
		if err := json.Unmarshal([]byte(rawValue), &tempValue); err == nil {
			value = tempValue
		}
	}

	if value == "" || value == nil {
		switch key {
		case "inactive":
			value = []string{}
		default:
			value = map[string]interface{}{}
		}
	}

	return value, nil
}

func fetchRedisHashValue(key, domainIDStr string, client *redis.Client, ctx context.Context) interface{} {
	rawValue, err := client.HGet(ctx, domainIDStr, key).Result()
	if err != nil {
		log.Printf("Error retrieving Redis hash for key %s: %v", key, err)
		return nil
	}
	return rawValue
}

func getShortcuts(domainIDStr string, userID float64, client *redis.Client, ctx context.Context) (map[string]interface{}, error) {
	rawValue, err := client.HGet(ctx, domainIDStr, "shortcuts").Result()
	if err != nil || rawValue == "" {
		return nil, err
	}

	shortcuts := make(map[string]interface{})
	if err := json.Unmarshal([]byte(rawValue), &shortcuts); err != nil {
		return nil, fmt.Errorf("error unmarshalling shortcuts: %v", err)
	}

	localSC, err := client.Get(ctx, fmt.Sprintf("amsc_%v_%v", domainIDStr, userID)).Result()
	if err == nil && localSC != "" {
		localMap := make(map[string]interface{})
		if err := json.Unmarshal([]byte(localSC), &localMap); err == nil {
			for k, v := range localMap {
				shortcuts[k] = v
			}
		}
	}

	return shortcuts, nil
}

func getSynonyms(domainID string, client *redis.Client, ctx context.Context, stx *Storage) (map[string]interface{}, error) {
	keywordSynonymsKey := "askme_keyword_synonyms"
	askmeKeywordSynonymsStr, err := client.Get(ctx, keywordSynonymsKey).Result()

	var synonymsGlobal interface{}
	if err == redis.Nil {

		synonymsGlobal, err = stx.GetKeywordSynonyms()
		if err != nil {
			return nil, fmt.Errorf("error getting keyword synonyms: %v", err)
		}

		if askmeKeywordSynonyms, ok := synonymsGlobal.(map[string]map[string]interface{}); ok {
			askmeKeywordSynonymsJSON, _ := json.Marshal(askmeKeywordSynonyms)
			client.Set(ctx, keywordSynonymsKey, askmeKeywordSynonymsJSON, 0)
		} else {
			return nil, fmt.Errorf("unexpected type for keyword synonyms")
		}
	} else if err == nil {

		if err := json.Unmarshal([]byte(askmeKeywordSynonymsStr), &synonymsGlobal); err != nil {
			return nil, fmt.Errorf("error unmarshalling Redis synonyms: %v", err)
		}
	}

	sysGlobal, ok := synonymsGlobal.(map[string]interface{})
	if !ok {
		if synonymsGlobalMap, ok := synonymsGlobal.(map[string]map[string]interface{}); ok {
			sysGlobal = make(map[string]interface{})
			for key, subMap := range synonymsGlobalMap {
				sysGlobal[key] = subMap
			}
		} else {
			return nil, fmt.Errorf("unexpected type for synonymsGlobal")
		}
	}

	userSynonyms, err := client.HGet(ctx, domainID, "synonyms").Result()
	if err != nil {
		log.Println("User Synonym not found on redis %w", err)
	}

	var synonyms map[string]interface{}
	if err := json.Unmarshal([]byte(userSynonyms), &synonyms); err != nil {
		log.Println("User unmarshal user sysnonym not found on redis %w", err)
	}

	if userMap, ok := sysGlobal["user"].(map[string]interface{}); ok {
		if userSynonymMap, ok := synonyms["user"].(map[string]interface{}); ok {
			for k, v := range userSynonymMap {
				userMap[k] = v
			}
		} else {
			log.Println("Error: 'user' key in synonyms is not a map[string]interface{}")
		}
	} else {
		sysGlobal["user"] = synonyms["user"]
	}

	return sysGlobal, nil
}

func getWords(domainID string, client *redis.Client, ctx context.Context) (map[string][]map[string]interface{}, error) {
	words := make(map[string][]map[string]interface{})
	var wg sync.WaitGroup
	var mu sync.Mutex

	keys := []string{"filters", "metadata", "askme_keyword_names"}
	for _, key := range keys {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()

			value, err := client.Get(ctx, fmt.Sprintf("%s<<-->>words<<-->>%s", domainID, key)).Result()
			if err == redis.Nil {
				return
			}

			var data []map[string]interface{}
			if err := json.Unmarshal([]byte(value), &data); err == nil {
				mu.Lock()
				for _, item := range data {
					if word, ok := item["word"].(string); ok {
						upperWord := strings.ToUpper(word)
						words[upperWord] = append(words[upperWord], item)
					}
				}
				mu.Unlock()
			}
		}(key)
	}

	wg.Wait()
	return words, nil
}
