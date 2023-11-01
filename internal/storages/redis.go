package storages

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Demacr/otus-hl-socialnetwork/internal/config"
	"github.com/Demacr/otus-hl-socialnetwork/internal/domain"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

const (
	MIN_IDLE_CONNECTIONS int = 200
	POOL_SIZE            int = 12000
	POOL_TIMEOUT         int = 240
)

type redisCache struct {
	redisClient *redis.Client
}

func NewRedisCache(cfg *config.RedisConfig) CacheRepository {
	return &redisCache{
		redisClient: redis.NewClient(&redis.Options{
			Addr:         cfg.Host,
			MinIdleConns: MIN_IDLE_CONNECTIONS,
			PoolSize:     POOL_SIZE,
			PoolTimeout:  time.Duration(POOL_TIMEOUT) * time.Second,
			Password:     cfg.Password,
			DB:           cfg.Database,
		}),
	}
}

func (rc *redisCache) GetPost(post_id int) (*domain.Post, error) {
	postBytes, err := rc.redisClient.Get(context.TODO(), postKey(post_id)).Bytes()
	if errors.Is(err, redis.Nil) {
		return &domain.Post{}, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "RedisCache.GetPost.Get")
	}
	post := &domain.Post{}
	if err = json.Unmarshal(postBytes, post); err != nil {
		return nil, errors.Wrap(err, "RedisCache.GetPost.json.Unmarshal")
	}

	return post, nil
}

func (rc *redisCache) SetPost(post *domain.Post) error {
	postBytes, err := json.Marshal(post)
	if err != nil {
		return errors.Wrap(err, "RedisCache.SetPost.json.Marshal")
	}

	if err = rc.redisClient.Set(context.TODO(), postKey(post.ID), postBytes, 0).Err(); err != nil {
		return errors.Wrap(err, "RedisCache.SetPost.Set")
	}

	return nil
}

func (rc *redisCache) DeletePost(post_id int) error {
	if err := rc.redisClient.Del(context.TODO(), postKey(post_id)).Err(); err != nil {
		return errors.Wrap(err, "RedisCache.DeletePost.Del")
	}

	return nil
}

func (rc *redisCache) AddToFeed(profile_id int, post_id ...int) error {
	if len(post_id) == 0 {
		return nil
	}

	postsId := make([]string, len(post_id))
	for index, post := range post_id {
		postsId[index] = strconv.Itoa(post)
	}

	if err := rc.redisClient.LPush(context.TODO(), friendFeedKey(profile_id), postsId).Err(); err != nil {
		return errors.Wrap(err, "RedisCache.AddToFeed.LPush")
	}

	if err := rc.redisClient.LTrim(context.TODO(), friendFeedKey(profile_id), 0, 999).Err(); err != nil {
		return errors.Wrap(err, "RedisCache.AddToFeed.LTrim")
	}

	return nil
}

func (rc *redisCache) RebuildFeed(profileId int, post_id ...int) error {
	if len(post_id) == 0 {
		return nil
	}

	postsId := make([]string, len(post_id))
	for index, post := range post_id {
		postsId[index] = strconv.Itoa(post)
	}

	pipe := rc.redisClient.TxPipeline()
	pipe.Del(context.TODO(), friendFeedKey(profileId))
	pipe.LPush(context.TODO(), friendFeedKey(profileId), postsId)
	_, err := pipe.Exec(context.TODO())
	if err != nil {
		return errors.Wrap(err, "RedisCache.RebuildFeed.Pipe.Exec")
	}

	return nil
}

func (rc *redisCache) GetFeed(profileId int) (result []int, err error) {
	var resStrings []string
	if resStrings, err = rc.redisClient.LRange(context.TODO(), friendFeedKey(profileId), 0, 999).Result(); err != nil {
		return nil, errors.Wrap(err, "RedisCache.GetFeed.LRange")
	}

	result = make([]int, len(resStrings))
	for index, resString := range resStrings {
		result[index], err = strconv.Atoi(resString)
		if err != nil {
			log.Println(errors.Wrap(err, "RedisCache.GetFeed.ConvertToInt"))
		}
	}

	return result, nil
}

func postKey(id int) string {
	return fmt.Sprintf("post:%d", id)
}

func friendFeedKey(id int) string {
	return fmt.Sprintf("friend:feed:%d", id)
}
