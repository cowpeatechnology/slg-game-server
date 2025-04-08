package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"slg-game-server/internal/network"

	"github.com/anthdm/hollywood/actor"
)

func main() {
	// 创建context用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建Hollywood actor系统
	config := actor.NewEngineConfig()
	engine, err := actor.NewEngine(config)
	if err != nil {
		log.Fatal("Failed to create engine:", err)
	}

	// 创建WebSocket Hub
	hub := network.NewHub(engine)
	go hub.Run()

	// 设置WebSocket路由
	http.HandleFunc("/ws", hub.HandleWebSocket)

	// 启动HTTP服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil, // 使用默认的DefaultServeMux
	}

	go func() {
		log.Println("Starting WebSocket server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server error:", err)
		}
	}()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	select {
	case <-sigChan:
		log.Println("Received shutdown signal")
	case <-ctx.Done():
		log.Println("Context cancelled")
	}

	// 优雅关闭
	log.Println("Shutting down server...")

	// 关闭HTTP服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}
}
