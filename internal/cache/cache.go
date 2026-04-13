package cache

import (
	"context"
	"encoding/json"

	"github.com/Aryan-Gupta4460/letstalk/internal/models"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func CacheMessage(rdb *redis.Client, room string, msg models.Message) error {
	key := "chat:room:" + room

	data, _ := json.Marshal(msg)

	// push to list
	err := rdb.LPush(ctx, key, data).Err()
	if err != nil {
		return err
	}

	// keep only last 20 messages
	rdb.LTrim(ctx, key, 0, 19)

	return nil
}

func GetCachedMessages(rdb *redis.Client, room string) ([][]byte, error) {
	key := "chat:room:" + room

	results, err := rdb.LRange(ctx, key, 0, 19).Result()
	if err != nil {
		return nil, err
	}

	var messages [][]byte
	for _, r := range results {
		messages = append(messages, []byte(r))
	}

	return messages, nil
}
