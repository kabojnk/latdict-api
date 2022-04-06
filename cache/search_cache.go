package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kabojnk/latdict-api/types"
	"log"
	"time"
)

type SearchCache struct {
	RedisClient     *redis.Client
	CacheTTLSeconds int
}

func (searchCache *SearchCache) IsEnabled() bool {
	redisConfig := RedisConfig{}
	redisConfig.Init()
	return redisConfig.ShouldDisableCache == false
}

func (searchCache *SearchCache) Open() {
	redisConfig := RedisConfig{}
	redisConfig.Init()
	searchCache.CacheTTLSeconds = redisConfig.CacheTTLSeconds
	searchCache.RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
}

func (searchCache *SearchCache) Close() {
	err := searchCache.RedisClient.Close()
	if err != nil {
		return
	}
	searchCache.RedisClient = nil
}

// GetSearchCacheKey Returns a formatted search cache key, used to index our cache based on certain criteria.
// @TODO it may be a good idea to revisit this caching strategy if pagination becomes more configurable to the public,
//       as the extra page number and page size dimensions might cause the cache to balloon out of control.
func GetSearchCacheKey(language string, term string, doesRequireExactMatch bool, pageNum int, pageSize int) string {
	exactMatchKey := "nonExact"
	if doesRequireExactMatch {
		exactMatchKey = "exact"
	}
	return fmt.Sprintf("%s:%s:%s:%d:%d", language, term, exactMatchKey, pageNum, pageSize)
}

func (searchCache *SearchCache) SaveCache(language string, term string, doesRequireExactMatch bool, pageNum int, pageSize int, v any) {
	jsonData, err := json.Marshal(v)
	ctx := context.Background()
	key := GetSearchCacheKey(language, term, doesRequireExactMatch, pageNum, pageSize)
	fmt.Printf("Saving cache entry for key: %s", key)
	err = searchCache.RedisClient.Set(ctx, key, jsonData, time.Duration(searchCache.CacheTTLSeconds)*time.Second).Err()
	if err != nil {
		log.Fatalf("Unable to save cache: %v\n", err)
	}
}

func (searchCache *SearchCache) GetCache(language string, term string, doesRequireExactMatch bool, pageNum int, pageSize int) (types.EntriesResponse, error) {
	ctx := context.Background()
	key := GetSearchCacheKey(language, term, doesRequireExactMatch, pageNum, pageSize)
	fmt.Printf("Lookupg up cache entry for key: %s", key)
	var entriesResponse types.EntriesResponse
	data, err := searchCache.RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return types.EntriesResponse{}, err
	}
	err = json.Unmarshal(data, &entriesResponse)
	fmt.Printf("FOUND cache entry for key: %s", key)
	if err != nil {
		return types.EntriesResponse{}, err
	}
	return entriesResponse, nil
}
