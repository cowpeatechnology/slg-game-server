# SLG Game Server

基于 Hollywood Actor 框架的 SLG 游戏服务器，使用 WebSocket 和 Protocol Buffers 实现实时通信。

## 项目结构

```
.
├── cmd/                    # 主程序入口
│   └── server/            # 服务器入口
│       └── main.go        # 主程序
├── internal/              # 内部包
│   ├── network/          # 网络相关
│   │   └── websocket.go  # WebSocket 处理
├── proto/                 # Protocol Buffers 定义
│   ├── message.proto     # 消息协议定义
│   └── message.pb.go     # 生成的 Go 代码
└── test/                 # 测试相关
    └── test.html         # WebSocket 客户端测试页面
```

## 技术栈

- Go 1.21+
- Hollywood Actor Framework
- WebSocket
- Protocol Buffers
- HTML/JavaScript (测试客户端)

## 消息协议

使用 Protocol Buffers 定义消息格式，支持以下消息类型：

1. 心跳消息 (HEARTBEAT)
   - 用于保持连接活跃
   - 包含时间戳字段

2. 文本消息 (TEXT)
   - 用于发送文本内容
   - 包含消息内容和时间戳

3. 错误消息 (ERROR)
   - 用于错误处理
   - 包含错误码和错误信息

## WebSocket 通信

### 服务器端

- 端口：8080
- 路径：/ws
- 支持多客户端连接
- 实现消息广播机制

### 客户端

- 使用原生 WebSocket API
- 集成 protobuf.js 处理消息
- 支持自动重连
- 实现消息发送和接收展示

## 开发环境设置

1. 安装依赖：
```bash
go mod download
```

2. 安装 Protocol Buffers 编译器：
```bash
# macOS
brew install protobuf
# 或其他平台对应的安装方式
```

3. 安装 Go Protocol Buffers 插件：
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

## 构建和运行

1. 生成 Protocol Buffers 代码：
```bash
protoc --go_out=. --go_opt=paths=source_relative proto/message.proto
```

2. 运行服务器：
```bash
go run cmd/server/main.go
```

3. 运行测试客户端：
   - 使用任意 Web 服务器托管 test/test.html
   - 或直接在浏览器中打开 test.html

## 测试

1. 打开测试客户端页面
2. 点击"连接"按钮连接到服务器
3. 可以发送文本消息或观察心跳
4. 打开多个客户端测试消息广播功能

## 注意事项

1. 生产环境部署时需要：
   - 配置适当的 CORS 策略
   - 添加安全验证机制
   - 实现完整的错误处理
   - 添加日志记录

2. 开发建议：
   - 遵循 Proto 文件定义的消息格式
   - 注意处理断线重连
   - 实现适当的错误处理机制

## 贡献

欢迎提交 Issue 和 Pull Request。

## 许可证

[MIT License](LICENSE) 