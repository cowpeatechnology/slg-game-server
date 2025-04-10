# SLG Game Server

基于 Actor 模型的简单游戏服务器框架，使用 Go 语言和 Hollywood Actor 框架实现。

## 设计原则

1. **单点写入原则**
   - 数据只能由特定 Actor 修改
   - 其他 Actor 通过消息请求修改
   - 示例：
     ```go
     // GameActor 负责管理玩家数据
     type GameActor struct {
         engine     *actor.Engine
         players    map[string]*pb.PlayerData
         combatPID  *actor.PID
     }

     // 其他 Actor 通过消息请求修改数据
     engine.Send(gameActor, &pb.GameMessage{
         Type: "update_player",
         Payload: playerData,
     })
     ```

2. **简单的消息流转**
   ```
   Client -> WebSocket -> GatewayActor -> GameActor/CombatActor
   ```
   - WebSocket 连接在基础设施层处理
   - GatewayActor 只负责消息路由
   - GameActor 和 CombatActor 处理具体业务逻辑

3. **数据存储**
   - 直接使用 Redis 存储
   - 简单的 CRUD 操作
   - 无需复杂的事务处理

## 项目结构

```
.
├── cmd/
│   └── server/
│       └── main.go           # 服务器入口，初始化组件
├── config/
│   └── config.json          # 基础配置（服务器、Redis）
├── internal/
│   ├── config/
│   │   └── config.go        # 配置加载和管理
│   ├── game/
│   │   ├── game_actor.go    # 游戏核心逻辑
│   │   └── combat_actor.go  # 战斗系统
│   ├── gateway/
│   │   └── gateway_actor.go # 消息路由和WebSocket连接管理
│   └── storage/
│       ├── redis.go         # Redis 存储实现
│       ├── storage.go       # 存储接口定义
│       └── storage_actor.go # 存储 Actor
├── proto/
│   ├── message.proto        # 消息协议定义
│   └── message.pb.go        # 生成的代码
└── test/
    └── client.html          # 测试客户端
```

## 核心组件

### 1. Gateway Actor
```go
// 简单的消息路由
type GatewayActor struct {
    engine    *actor.Engine
    clients   sync.Map
    gameActor *actor.PID
}

// 转发消息到对应的 Actor
func (a *GatewayActor) handleMessage(msg *pb.GameMessage) {
    a.engine.Send(a.gameActor, msg)
}
```

### 2. Game Actor
```go
// 游戏核心逻辑
type GameActor struct {
    engine     *actor.Engine
    players    map[string]*pb.PlayerData
    combatPID  *actor.PID
}

// 处理玩家数据和游戏逻辑
func (a *GameActor) handlePlayerJoin(msg *pb.GameMessage) {
    // 处理玩家加入逻辑
}
```

### 3. Combat Actor
```go
// 战斗系统
type CombatActor struct {
    engine    *actor.Engine
    gamePID   *actor.PID
    battles   map[string]*pb.BattleResult
}

// 处理战斗逻辑
func (a *CombatActor) handleBattle(msg *pb.GameMessage) {
    // 处理战斗逻辑
}
```

## 数据结构

```protobuf
// 统一的玩家数据结构
message PlayerData {
    string id = 1;
    string name = 2;
    int32 level = 3;
    int32 hp = 4;
    int32 attack = 5;
    int32 defense = 6;
}

// 游戏消息
message GameMessage {
    string type = 1;
    bytes payload = 2;
    string id = 3;
}
```

## 开发环境要求

- Go 1.21+
- Redis 6.0+
- hollywood v1.0.5
- protoc v4.25.3+

## 快速开始

1. 安装依赖
```bash
go mod tidy
```

2. 生成 protobuf 代码
```bash
protoc --go_out=. --go_opt=paths=source_relative proto/message.proto
```

3. 启动服务器
```bash
go run cmd/server/main.go
```

4. 测试
- 打开 `test/client.html` 进行测试
- 使用 WebSocket 连接到服务器
- 发送消息测试功能

## 消息流程示例

1. **玩家加入游戏**
```
Client -> GatewayActor -> GameActor
- GatewayActor 记录客户端连接
- GameActor 创建/加载玩家数据
```

2. **战斗请求**
```
Client -> GatewayActor -> GameActor -> CombatActor
- GameActor 验证玩家状态
- CombatActor 处理战斗逻辑
- 结果通过相同路径返回
```

## 注意事项

1. **Actor 通信**
   - 始终通过 PID 发送消息
   - 不要在 Actor 之间直接调用方法
   - 保持消息处理的异步性

2. **数据处理**
   - GameActor 是玩家数据的唯一修改者
   - 使用 Redis 进行简单的数据持久化
   - 避免复杂的数据查询操作

3. **错误处理**
   - 记录关键错误日志
   - 保持系统稳定性
   - 优雅处理连接断开

## 调试建议

1. 使用日志跟踪消息流转
2. 监控 Redis 连接状态
3. 使用测试客户端验证功能

## 部署

1. 准备环境
   - 配置 Redis
   - 设置配置文件

2. 构建
```bash
go build -o server cmd/server/main.go
```

3. 运行
```bash
./server
``` 