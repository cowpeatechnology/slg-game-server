package game

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/anthdm/hollywood/actor"
	pb "github.com/cowpeatechnology/slg-game-server/proto"
)

// GameActor handles game logic
type GameActor struct {
	engine     *actor.Engine
	players    map[string]*pb.PlayerData
	gatewayPID *actor.PID
	combatPID  *actor.PID
}

// NewGameActor creates a new Game Actor
func NewGameActor() actor.Producer {
	return func() actor.Receiver {
		return &GameActor{
			players: make(map[string]*pb.PlayerData),
		}
	}
}

// Receive handles incoming messages
func (a *GameActor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		log.Println("[GameActor] Started")
		a.engine = ctx.Engine()

	case actor.Stopped:
		log.Println("[GameActor] Stopped")

	case *actor.PID:
		// log.Printf("[GameActor] Received PID: %v", msg)
		// log.Printf("[GameActor] Received msg.ID: %v", msg.ID)
		if msg.ID != "" {
			if strings.Contains(msg.ID, "combat/") {
				a.combatPID = msg
				log.Printf("[GameActor] Received CombatActor PID: %v", msg)
			}
			if strings.Contains(msg.ID, "gateway/") {
				a.gatewayPID = msg
				log.Printf("[GameActor] Received GatewayActor PID: %v", msg)
			}
		}

	case *pb.GameMessage:
		a.handleGameMessage(ctx, msg)
	}
}

// handleGameMessage processes game-related messages
func (a *GameActor) handleGameMessage(ctx *actor.Context, msg *pb.GameMessage) {
	if err := a.validateMessage(msg); err != nil {
		log.Printf("[GameActor] 消息验证失败: %v", err)
		a.sendError(ctx, msg.Id, msg.Type, err)
		return
	}

	switch msg.Type {
	case "player_join":
		a.handlePlayerJoin(ctx, msg)
	case "chat":
		a.handleChat(ctx, msg)
	case "battle_request":
		if a.combatPID != nil {
			ctx.Engine().Send(a.combatPID, msg)
		} else {
			log.Printf("[GameActor] CombatActor PID not available")
			a.sendError(ctx, msg.Id, msg.Type, fmt.Errorf("combat service not available"))
		}
	case "battle_result":
		// 解析战斗结果，获取参与战斗的玩家ID
		var battleResult struct {
			WinnerID    string `json:"winner_id"`
			LoserID     string `json:"loser_id"`
			DamageDealt int32  `json:"damage_dealt"`
		}
		if err := json.Unmarshal(msg.Payload, &battleResult); err != nil {
			log.Printf("[GameActor] Failed to unmarshal battle result: %v", err)
			return
		}

		// 确保两个玩家都在线
		_, winnerExists := a.players[battleResult.WinnerID]
		_, loserExists := a.players[battleResult.LoserID]

		if !winnerExists || !loserExists {
			log.Printf("[GameActor] One or both players not found: winner=%s, loser=%s",
				battleResult.WinnerID, battleResult.LoserID)
			return
		}

		// 分别发送战斗结果给胜利者和失败者
		if a.gatewayPID != nil {
			// 发送给胜利者
			winnerMsg := &pb.GameMessage{
				Type:    "battle_result",
				Id:      battleResult.WinnerID,
				Payload: msg.Payload,
			}
			ctx.Engine().Send(a.gatewayPID, winnerMsg)
			log.Printf("[GameActor] Sent battle result to winner: %s", battleResult.WinnerID)

			// 发送给失败者
			loserMsg := &pb.GameMessage{
				Type:    "battle_result",
				Id:      battleResult.LoserID,
				Payload: msg.Payload,
			}
			ctx.Engine().Send(a.gatewayPID, loserMsg)
			log.Printf("[GameActor] Sent battle result to loser: %s", battleResult.LoserID)
		}
	default:
		log.Printf("[GameActor] Unknown message type: %s", msg.Type)
	}
}

// handlePlayerJoin handles player join requests
func (a *GameActor) handlePlayerJoin(ctx *actor.Context, msg *pb.GameMessage) {
	// 生成玩家ID
	playerID := fmt.Sprintf("player_%d", len(a.players)+1)

	// 创建新玩家数据
	player := &pb.PlayerData{
		Id:      playerID,
		Name:    fmt.Sprintf("Player_%d", len(a.players)+1),
		Level:   1,
		Hp:      100,
		Attack:  10,
		Defense: 5,
	}

	// 保存玩家数据
	a.players[playerID] = player

	// 创建响应消息
	response := &pb.GameMessage{
		Type: "player_join_response",
		Id:   playerID,
		Payload: []byte(fmt.Sprintf(`{"id":"%s","name":"%s"}`,
			player.Id, player.Name)),
	}

	// 发送响应
	if a.gatewayPID != nil {
		ctx.Engine().Send(a.gatewayPID, response)
	}

	log.Printf("[GameActor] New player joined: %s", playerID)
}

// handleChat processes chat messages
func (a *GameActor) handleChat(ctx *actor.Context, msg *pb.GameMessage) {
	// 确保发送者是已登录的玩家
	sender, exists := a.players[msg.Id]
	if !exists {
		a.sendError(ctx, msg.Id, msg.Type, fmt.Errorf("player not found"))
		return
	}

	// 创建聊天响应消息
	response := &pb.GameMessage{
		Type: "chat_response",
		Id:   msg.Id,
		Payload: []byte(fmt.Sprintf(`{
			"sender": "%s",
			"sender_name": "%s",
			"content": "%s",
			"timestamp": %d
		}`, sender.Id, sender.Name, string(msg.Payload), time.Now().Unix())),
	}

	// 广播消息给所有在线玩家
	if a.gatewayPID != nil {
		ctx.Engine().Send(a.gatewayPID, response)
	}
}

// validateMessage 验证消息的有效性
func (a *GameActor) validateMessage(msg *pb.GameMessage) error {
	if msg == nil {
		return fmt.Errorf("消息为空")
	}
	if msg.Type == "" {
		return fmt.Errorf("消息类型为空")
	}
	return nil
}

// sendError 发送错误消息给客户端
func (a *GameActor) sendError(ctx *actor.Context, id string, msgType string, err error) {
	if a.gatewayPID == nil {
		log.Printf("[GameActor] 无法发送错误消息: gatewayPID为空")
		return
	}

	errMsg := &pb.GameMessage{
		Type:    "error",
		Id:      id,
		Payload: []byte(err.Error()),
	}

	ctx.Engine().Send(a.gatewayPID, errMsg)
}
