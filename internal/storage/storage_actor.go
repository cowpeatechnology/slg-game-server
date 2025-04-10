package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/anthdm/hollywood/actor"
	pb "github.com/cowpeatechnology/slg-game-server/proto"
	"github.com/go-redis/redis/v8"
	protobuf "google.golang.org/protobuf/proto"
)

// StorageActor 处理数据存储
type StorageActor struct {
	engine *actor.Engine
	redis  *redis.Client
}

// NewStorageActor 创建 Storage Actor
func NewStorageActor(redisClient *redis.Client) actor.Receiver {
	return &StorageActor{
		redis: redisClient,
	}
}

// Receive 处理接收到的消息
func (a *StorageActor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		log.Println("StorageActor started")
		a.engine = ctx.Engine()

	case *StorageRequestMessage:
		// 处理存储请求
		a.handleStorageRequest(ctx, msg)
	}
}

// StorageRequestMessage 存储请求消息
type StorageRequestMessage struct {
	Type string      // 操作类型：get_player, save_player, delete_player
	Key  string      // 玩家ID
	Data interface{} // 玩家数据
}

// StorageResponseMessage 存储响应消息
type StorageResponseMessage struct {
	Type  string      // 操作类型
	Data  interface{} // 返回的数据
	Error error       // 错误信息
}

// 处理存储请求
func (a *StorageActor) handleStorageRequest(ctx *actor.Context, msg *StorageRequestMessage) {
	var response StorageResponseMessage
	response.Type = msg.Type

	switch msg.Type {
	case "get_player":
		// 获取玩家数据
		data, err := a.redis.Get(context.Background(), a.getPlayerKey(msg.Key)).Bytes()
		if err != nil {
			if err == redis.Nil {
				response.Error = fmt.Errorf("玩家不存在: %s", msg.Key)
			} else {
				response.Error = fmt.Errorf("获取玩家数据失败: %v", err)
			}
		} else {
			var player pb.PlayerData
			if err := protobuf.Unmarshal(data, &player); err != nil {
				response.Error = fmt.Errorf("解析玩家数据失败: %v", err)
			} else {
				response.Data = &player
			}
		}

	case "save_player":
		// 保存玩家数据
		if player, ok := msg.Data.(*pb.PlayerData); ok {
			data, err := protobuf.Marshal(player)
			if err != nil {
				response.Error = fmt.Errorf("序列化玩家数据失败: %v", err)
			} else {
				err = a.redis.Set(context.Background(), a.getPlayerKey(msg.Key), data, 24*time.Hour).Err()
				if err != nil {
					response.Error = fmt.Errorf("保存玩家数据失败: %v", err)
				}
			}
		} else {
			response.Error = fmt.Errorf("无效的玩家数据类型")
		}

	case "delete_player":
		// 删除玩家数据
		err := a.redis.Del(context.Background(), a.getPlayerKey(msg.Key)).Err()
		if err != nil {
			response.Error = fmt.Errorf("删除玩家数据失败: %v", err)
		}

	default:
		response.Error = fmt.Errorf("未知的操作类型: %s", msg.Type)
	}

	// 发送响应
	if ctx.Sender() != nil {
		a.engine.Send(ctx.Sender(), &response)
	}

	// 记录操作日志
	if response.Error != nil {
		log.Printf("[StorageActor] 操作失败: type=%s, key=%s, error=%v", msg.Type, msg.Key, response.Error)
	} else {
		log.Printf("[StorageActor] 操作成功: type=%s, key=%s", msg.Type, msg.Key)
	}
}

// getPlayerKey generates a Redis key for a player
func (a *StorageActor) getPlayerKey(id string) string {
	return fmt.Sprintf("player:%s", id)
}
