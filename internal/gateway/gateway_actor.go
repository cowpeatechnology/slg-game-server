package gateway

import (
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/anthdm/hollywood/actor"
	pb "github.com/cowpeatechnology/slg-game-server/proto"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

// ConnectMessage is sent when a new client connects
type ConnectMessage struct {
	ClientID string
	Conn     *websocket.Conn
}

// GatewayActor handles WebSocket connections and message routing
type GatewayActor struct {
	engine    *actor.Engine
	clients   sync.Map // key: clientID, value: *websocket.Conn
	players   sync.Map // key: playerID, value: clientID
	gameActor *actor.PID
}

// NewGatewayActor creates a new Gateway Actor
func NewGatewayActor() actor.Producer {
	return func() actor.Receiver {
		return &GatewayActor{}
	}
}

// Receive handles incoming messages
func (a *GatewayActor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		log.Println("[GatewayActor] Started")
		a.engine = ctx.Engine()

	case actor.Stopped:
		log.Println("[GatewayActor] Stopped")
		// 清理所有连接
		a.clients.Range(func(key, value interface{}) bool {
			if conn, ok := value.(*websocket.Conn); ok {
				conn.Close()
			}
			return true
		})

	case *actor.PID:
		// log.Printf("[GatewayActor] Received PID: %v", msg)
		// log.Printf("[GatewayActor] Received msg.ID: %v", msg.ID)
		if msg.ID != "" && strings.Contains(msg.ID, "game/") {
			a.gameActor = msg
			log.Printf("[GatewayActor] Received GameActor PID: %v", msg)
		}

	case *ConnectMessage:
		// 存储连接
		a.clients.Store(msg.ClientID, msg.Conn)
		log.Printf("[GatewayActor] Client connected: %s", msg.ClientID)
		// 启动消息读取
		go a.readPump(ctx, msg.ClientID, msg.Conn)

	case []byte:
		// 处理从WebSocket接收到的原始消息
		a.handleWebSocketMessage(ctx, msg)

	case *pb.GameMessage:
		a.handleGameMessage(ctx, msg)
	}
}

// readPump reads messages from the WebSocket connection
func (a *GatewayActor) readPump(ctx *actor.Context, clientID string, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		a.clients.Delete(clientID)
		log.Printf("[GatewayActor] Client disconnected: %s", clientID)
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[GatewayActor] Read error: %v", err)
			return
		}

		// 直接发送原始数据给自身处理
		a.engine.Send(ctx.PID(), data)
	}
}

// handleWebSocketMessage processes messages received from WebSocket
func (a *GatewayActor) handleWebSocketMessage(ctx *actor.Context, data []byte) {
	if a.gameActor == nil {
		log.Printf("[GatewayActor] No GameActor available")
		return
	}

	// 尝试解析为GameMessage
	var gameMsg pb.GameMessage
	if err := proto.Unmarshal(data, &gameMsg); err != nil {
		log.Printf("[GatewayActor] Failed to unmarshal message: %v", err)
		return
	}

	// 转发消息到 GameActor
	a.engine.Send(a.gameActor, &gameMsg)
}

// handleGameMessage processes messages from GameActor
func (a *GatewayActor) handleGameMessage(ctx *actor.Context, msg *pb.GameMessage) {
	// 对于 player_join_response 消息，需要建立 playerID 到 clientID 的映射
	if msg.Type == "player_join_response" {
		// 解析响应消息获取玩家ID
		var response struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(msg.Payload, &response); err != nil {
			log.Printf("[GatewayActor] Failed to unmarshal player join response: %v", err)
			return
		}

		// 遍历所有客户端连接，找到未映射的连接
		var foundClientID string
		a.clients.Range(func(key, _ interface{}) bool {
			clientID := key.(string)
			// 检查这个clientID是否已经被映射到其他playerID
			var isMapped bool
			a.players.Range(func(_, value interface{}) bool {
				if value.(string) == clientID {
					isMapped = true
					return false
				}
				return true
			})
			if !isMapped {
				foundClientID = clientID
				return false
			}
			return true
		})

		if foundClientID != "" {
			// 建立 playerID 到 clientID 的映射
			a.players.Store(response.ID, foundClientID)
			log.Printf("[GatewayActor] Mapped player %s to client %s", response.ID, foundClientID)
		}
	}

	// 通过 playerID 查找 clientID
	if clientID, ok := a.players.Load(msg.Id); ok {
		// 通过 clientID 查找连接
		if conn, ok := a.clients.Load(clientID); ok {
			wsConn := conn.(*websocket.Conn)
			// 序列化消息
			data, err := proto.Marshal(msg)
			if err != nil {
				log.Printf("[GatewayActor] Error encoding message: %v", err)
				return
			}

			// 发送消息给客户端
			if err := wsConn.WriteMessage(websocket.BinaryMessage, data); err != nil {
				log.Printf("[GatewayActor] Error sending message to client %s: %v", clientID, err)
				// 连接可能已断开，清理映射
				a.clients.Delete(clientID)
				a.players.Delete(msg.Id)
			}
		}
	} else {
		log.Printf("[GatewayActor] Player not found: %s", msg.Id)
	}
}
