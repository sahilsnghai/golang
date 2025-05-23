package services

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
)

type Storage struct {
	Db         *sql.DB
	Parameters map[string]interface{}
	Client     *redis.Client
	Ctx        context.Context
	Headers    map[string]string
}

func (s *Storage) QueryWithTiming(query string, args ...interface{}) (*sql.Rows, error) {
	startTime := time.Now()
	rows, err := s.Db.Query(query, args...)
	elapsedTime := time.Since(startTime)
	log.Printf("\n\nQuery executed in: %v\nSQL: %s\nArgs: %v", elapsedTime, query, args)
	return rows, err
}

func (s *Storage) hPublish(domainIDStr, key string, resp any) {
	log.Printf("publishing %s for %s\n", key, domainIDStr)
	_resp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error marshalling data to JSON: %v", err)
	}

	err = s.Client.HSet(s.Ctx, domainIDStr, key, _resp).Err()
	if err != nil {
		log.Fatalf("Error setting %s in Redis: %v", key, err)
	}
	log.Printf("publishing %s is completed\n", key)
}

func (s *Storage) Publish(domainIDStr, key string, resp any) {
	log.Printf("publishing %s for %s\n", key, domainIDStr)
	_resp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error marshalling datta to JSON: %v", err)
	}

	err = s.Client.Set(s.Ctx, key, _resp, 0).Err()
	if err != nil {
		log.Fatalf("Error setting %s in Redis: %v", key, err)
	}
	log.Printf("publishing %s is completed\n", key)
}

func (s *Storage) calculateFilterLimit() int {
	askmeFilterLimits, ok := s.Parameters["askme_filter_limits"].(map[string]interface{})
	if !ok {
		log.Println("Error: askme_filter_limits is not a map[string]interface{}")
		return 1000
	}

	FILTER_MAX_LIMIT, _ := askmeFilterLimits["max_limit"].(int)
	FILTER_MIN_LIMIT, _ := askmeFilterLimits["min_limit"].(int)
	FILTER_PER_MB, _ := askmeFilterLimits["filters_per_mb"].(int)

	if FILTER_MAX_LIMIT == 0 {
		FILTER_MAX_LIMIT = 1000
	}
	if FILTER_MIN_LIMIT == 0 {
		FILTER_MIN_LIMIT = 1000
	}
	if FILTER_PER_MB == 0 {
		FILTER_PER_MB = 10
	}

	info, err := s.Client.Info(s.Ctx, "memory").Result()
	if err != nil {
		log.Printf("Error retrieving memory info from Redis: %v", err)
		return FILTER_MIN_LIMIT
	}

	var maxMemory, totalSystemMemory, usedMemory int64
	for _, line := range strings.Split(info, "\n") {
		switch {
		case strings.HasPrefix(line, "maxmemory:"):
			_, _ = fmt.Sscanf(line, "maxmemory:%d", &maxMemory)
		case strings.HasPrefix(line, "total_system_memory:"):
			_, _ = fmt.Sscanf(line, "total_system_memory:%d", &totalSystemMemory)
		case strings.HasPrefix(line, "used_memory:"):
			_, _ = fmt.Sscanf(line, "used_memory:%d", &usedMemory)
		}
	}

	totalMemory := maxMemory
	if maxMemory == 0 {
		totalMemory = totalSystemMemory
	}

	freeMemory := totalMemory - usedMemory
	log.Printf("FILTER_MAX_LIMIT: %d, FILTER_PER_MB * freeMemory: %.f,  freeMemory %d", FILTER_MAX_LIMIT, float64(FILTER_PER_MB)*float64(freeMemory), freeMemory)
	return int(math.Min(
		float64(FILTER_MAX_LIMIT),
		float64(FILTER_PER_MB)*float64(freeMemory),
	))
}

func cleanQuery(text string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.ToUpper(re.ReplaceAllString(strings.ReplaceAll(strings.ReplaceAll(text, ",", " "), "'", ""), " "))
}

func GetFromRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	askmeWord := make([]map[string]interface{}, 0)
	for rows.Next() {
		k, err := scanRows(rows)
		if err != nil {
			return nil, err
		}
		askmeWord = append(askmeWord, k)
	}
	return askmeWord, nil
}

func scanRows(rows *sql.Rows) (map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	result := make(map[string]interface{}, len(columns))

	for i := range values {
		values[i] = new(interface{})
	}
	if err := rows.Scan(values...); err != nil {
		return nil, err
	}
	for i, colName := range columns {
		v := *values[i].(*interface{})

		switch v := v.(type) {
		case []byte:
			if len(v) == 1 {
				result[colName] = v[0]
			} else {
				result[colName] = string(v)
			}
		case nil:
			result[colName] = nil
		default:
			result[colName] = v
		}
	}

	return result, nil
}

func (s *Storage) UpdateReqCtx(ctx *fasthttp.RequestCtx) {
	if s.Headers == nil {
		s.Headers = make(map[string]string)
	}
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		s.Headers[string(key)] = string(value)
	})

}

func (s *Storage) MakeRequest(method, url string, payload map[string]interface{}, headers map[string]string) (*any, error) {
	var jsonPayload []byte
	var err error

	if payload != nil {
		jsonPayload, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("error marshaling payload: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	maps.Copy(headers, s.Headers)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	log.Printf("All the values of req.headers %v", req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received non-OK status code %d", resp.StatusCode)
	}

	var responseBody any
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return &responseBody, nil
}
