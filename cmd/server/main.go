package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/cowpeatechnology/slg-game-server/internal/config"
	"github.com/cowpeatechnology/slg-game-server/internal/game"
	"github.com/cowpeatechnology/slg-game-server/internal/gateway"
	"github.com/gorilla/websocket"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize actor engine
	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatalf("Failed to create actor engine: %v", err)
	}

	// Initialize actors
	gameActor := engine.Spawn(game.NewGameActor(), "game")
	combatActor := engine.Spawn(game.NewCombatActor(), "combat")
	gatewayActor := engine.Spawn(gateway.NewGatewayActor(), "gateway")

	// 等待Actor完全启动
	time.Sleep(100 * time.Millisecond)
	log.Printf("Actor PIDs started - Game: %v, Combat: %v, Gateway: %v",
		gameActor, combatActor, gatewayActor)

	// 设置Actor之间的PID引用
	// 首先发送 GameActor 的 PID 给其他 Actor
	engine.Send(gatewayActor, gameActor) // Gateway 需要知道 Game 的 PID
	engine.Send(combatActor, gameActor)  // Combat 需要知道 Game 的 PID

	// 然后发送其他 Actor 的 PID 给 GameActor
	engine.Send(gameActor, gatewayActor) // Game 需要知道 Gateway 的 PID
	engine.Send(gameActor, combatActor)  // Game 需要知道 Combat 的 PID

	log.Printf("Actor PIDs exchanged - Game: %v, Combat: %v, Gateway: %v",
		gameActor, combatActor, gatewayActor)

	// Initialize HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	mux := http.NewServeMux()

	// Create WebSocket upgrader
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}

		clientID := fmt.Sprintf("%s-%s", r.RemoteAddr, conn.LocalAddr().String())
		engine.Send(gatewayActor, &gateway.ConnectMessage{
			ClientID: clientID,
			Conn:     conn,
		})
	})

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Start HTTP server
	log.Printf("Starting gateway service on %s", addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Graceful shutdown
	log.Println("Shutting down server...")
	server.Close()
}
