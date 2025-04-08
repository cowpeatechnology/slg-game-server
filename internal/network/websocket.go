package network

import (
	"log"
	"net/http"
	"sync"

	pb "slg-game-server/proto"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

// Client 表示一个WebSocket客户端连接
type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	hub    *Hub
	engine *actor.Engine
	mu     sync.Mutex
}

// Hub 管理所有活动的WebSocket连接
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	engine     *actor.Engine
	upgrader   websocket.Upgrader
}

// NewHub 创建一个新的Hub
func NewHub(engine *actor.Engine) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		engine:     engine,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许所有来源的连接，生产环境中应该更严格
			},
		},
	}
}

// Run 启动Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client connected. Total clients: %d", len(h.clients))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client disconnected. Total clients: %d", len(h.clients))
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// HandleWebSocket 处理WebSocket连接请求
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 256),
		hub:    h,
		engine: h.engine,
	}

	h.register <- client

	// 启动goroutine处理读写
	go client.writePump()
	go client.readPump()
}

// readPump 从WebSocket连接读取消息
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		// 解析protobuf消息
		var msg pb.Message
		if err := proto.Unmarshal(message, &msg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// 根据消息类型处理不同的消息
		log.Printf("Received message type: %v", msg.Type)

		switch msg.Type {
		case pb.MessageType_HEARTBEAT:
			// 解析心跳消息
			var heartbeat pb.Heartbeat
			if err := proto.Unmarshal(msg.Payload, &heartbeat); err != nil {
				log.Printf("Failed to unmarshal heartbeat: %v", err)
				continue
			}

			log.Printf("Received heartbeat with timestamp: %v", heartbeat.Timestamp)

			// 创建响应消息
			responseMsg := &pb.Message{
				Type:    pb.MessageType_HEARTBEAT,
				Payload: msg.Payload, // 回显相同的时间戳
			}

			// 序列化响应消息
			responseData, err := proto.Marshal(responseMsg)
			if err != nil {
				log.Printf("Failed to marshal response: %v", err)
				continue
			}

			// 发送响应
			c.send <- responseData

		case pb.MessageType_TEXT:
			// 解析文本消息
			var textMessage pb.TextMessage
			if err := proto.Unmarshal(msg.Payload, &textMessage); err != nil {
				log.Printf("Failed to unmarshal text message: %v", err)
				continue
			}

			log.Printf("Received text message: %v (timestamp: %v)", textMessage.Content, textMessage.Timestamp)

			// 创建响应消息（回显）
			responseMsg := &pb.Message{
				Type:    pb.MessageType_TEXT,
				Payload: msg.Payload,
			}

			// 序列化响应消息
			responseData, err := proto.Marshal(responseMsg)
			if err != nil {
				log.Printf("Failed to marshal response: %v", err)
				continue
			}

			// 发送响应给所有客户端（广播）
			c.hub.broadcast <- responseData

		default:
			log.Printf("Unknown message type: %v", msg.Type)
		}
	}
}

// writePump 向WebSocket连接写入消息
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.mu.Lock()
			if !ok {
				// 通道已关闭
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				c.mu.Unlock()
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				c.mu.Unlock()
				return
			}
			w.Write(message)

			// 添加队列中的其他消息
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				c.mu.Unlock()
				return
			}
			c.mu.Unlock()
		}
	}
}
