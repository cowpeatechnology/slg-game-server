package game

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"

	"github.com/anthdm/hollywood/actor"
	pb "github.com/cowpeatechnology/slg-game-server/proto"
)

// CombatActor handles battle logic
type CombatActor struct {
	engine     *actor.Engine
	combatData sync.Map   // 战斗数据，key为id
	gamePID    *actor.PID // 游戏Actor的引用
	battles    map[string]*pb.BattleResult
}

// NewCombatActor creates a new Combat Actor
func NewCombatActor() actor.Producer {
	return func() actor.Receiver {
		return &CombatActor{
			battles: make(map[string]*pb.BattleResult),
		}
	}
}

// Receive handles incoming messages
func (a *CombatActor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		log.Println("[CombatActor] Started")
		a.engine = ctx.Engine()

	case actor.Stopped:
		log.Println("[CombatActor] Stopped")

	case *actor.PID:
		// log.Printf("[CombatActor] Received PID: %v", msg)
		// log.Printf("[CombatActor] Received msg.ID: %v", msg.ID)
		if msg.ID != "" && strings.Contains(msg.ID, "game/") {
			a.gamePID = msg
			log.Printf("[CombatActor] Received GameActor PID: %v", msg)
		}

	case *pb.GameMessage:
		if err := a.validateMessage(msg); err != nil {
			log.Printf("[CombatActor] 消息验证失败: %v", err)
			a.sendError(ctx, msg.Id, msg.Type, err)
			return
		}
		log.Printf("[CombatActor] 收到消息: type=%s, id=%s", msg.Type, msg.Id)
		a.handleCombatMessage(ctx, msg)
	}
}

// validateMessage 验证消息的有效性
func (a *CombatActor) validateMessage(msg *pb.GameMessage) error {
	if msg == nil {
		return fmt.Errorf("消息为空")
	}
	if msg.Type == "" {
		return fmt.Errorf("消息类型为空")
	}
	if msg.Type != "player_list" && msg.Id == "" {
		return fmt.Errorf("玩家ID为空")
	}
	return nil
}

// sendError 发送错误消息给GameActor
func (a *CombatActor) sendError(ctx *actor.Context, id string, msgType string, err error) {
	if a.gamePID == nil {
		log.Printf("[CombatActor] 无法发送错误消息: gamePID为空")
		return
	}

	errMsg := &pb.GameMessage{
		Type:    msgType + "_error",
		Id:      id,
		Payload: []byte(err.Error()),
	}

	ctx.Engine().Send(a.gamePID, errMsg)
}

// handleCombatMessage processes combat-related messages
func (a *CombatActor) handleCombatMessage(ctx *actor.Context, msg *pb.GameMessage) {
	switch msg.Type {
	case "battle_request":
		var battleReq struct {
			AttackerID string `json:"attacker_id"`
			DefenderID string `json:"defender_id"`
		}
		if err := json.Unmarshal(msg.Payload, &battleReq); err != nil {
			log.Printf("[CombatActor] Failed to unmarshal battle request: %v", err)
			a.sendError(ctx, msg.Id, msg.Type, fmt.Errorf("invalid battle request format"))
			return
		}
		log.Printf("[CombatActor] Received battle request: attacker=%s, defender=%s",
			battleReq.AttackerID, battleReq.DefenderID)

		// Simple battle logic: random winner and damage
		var result pb.BattleResult
		if rand.Float32() > 0.5 {
			result.WinnerId = battleReq.AttackerID
			result.LoserId = battleReq.DefenderID
		} else {
			result.WinnerId = battleReq.DefenderID
			result.LoserId = battleReq.AttackerID
		}
		result.DamageDealt = int32(rand.Intn(50) + 10)

		// Serialize battle result
		resultBytes, err := json.Marshal(&result)
		if err != nil {
			log.Printf("[CombatActor] Failed to marshal battle result: %v", err)
			return
		}

		// Send battle result to GameActor
		if a.gamePID != nil {
			battleResultMsg := &pb.GameMessage{
				Type:    "battle_result",
				Id:      msg.Id, // 保持原始消息的ID
				Payload: resultBytes,
			}
			log.Printf("[CombatActor] Sending battle result: winner=%s, loser=%s, damage=%d",
				result.WinnerId, result.LoserId, result.DamageDealt)
			ctx.Engine().Send(a.gamePID, battleResultMsg)
		} else {
			log.Printf("[CombatActor] Cannot send battle result: GameActor PID not available")
		}

	default:
		err := fmt.Errorf("未知的消息类型: %s", msg.Type)
		a.sendError(ctx, msg.Id, msg.Type, err)
	}
}

// getCombatData 获取玩家的战斗数据
func (a *CombatActor) getCombatData(id string) (*pb.PlayerData, bool) {
	if data, ok := a.combatData.Load(id); ok {
		if playerData, ok := data.(*pb.PlayerData); ok {
			return playerData, true
		}
	}
	return nil, false
}
