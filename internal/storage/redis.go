package storage

import (
	"context"
	"fmt"

	pb "github.com/cowpeatechnology/slg-game-server/proto"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

// RedisConfig contains Redis connection configuration
type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

// RedisStorage implements the Storage interface using Redis
type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisStorageFactory creates a new Redis storage factory
func NewRedisStorageFactory(config RedisConfig) StorageFactory {
	return &redisStorageFactory{config: config}
}

type redisStorageFactory struct {
	config RedisConfig
}

func (f *redisStorageFactory) CreateStorage() (Storage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     f.config.Address,
		Password: f.config.Password,
		DB:       f.config.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisStorage{
		client: client,
		ctx:    ctx,
	}, nil
}

// GetPlayer retrieves player data from Redis
func (s *RedisStorage) GetPlayer(id string) (*pb.PlayerData, error) {
	data, err := s.client.Get(s.ctx, fmt.Sprintf("player:%s", id)).Bytes()
	if err != nil {
		return nil, err
	}

	var player pb.PlayerData
	if err := proto.Unmarshal(data, &player); err != nil {
		return nil, err
	}

	return &player, nil
}

func (s *RedisStorage) SavePlayer(player *pb.PlayerData) error {
	data, err := proto.Marshal(player)
	if err != nil {
		return err
	}

	return s.client.Set(s.ctx, fmt.Sprintf("player:%s", player.Id), data, 0).Err()
}

func (s *RedisStorage) UpdatePlayer(player *pb.PlayerData) error {
	return s.SavePlayer(player)
}

// DeletePlayer deletes player data from Redis
func (s *RedisStorage) DeletePlayer(id string) error {
	return s.client.Del(s.ctx, fmt.Sprintf("player:%s", id)).Err()
}

func (s *RedisStorage) Close() error {
	return s.client.Close()
}
