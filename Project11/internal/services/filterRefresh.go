package services

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func (s *Storage) FilterRefresh(userId, domainId, orgId int, ctx context.Context, client *redis.Client) (map[string]interface{}, error) {
	return nil, nil
}
