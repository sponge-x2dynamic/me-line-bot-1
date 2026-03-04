package main

import (
	"log"
	"me-bot/internal/config"
	"me-bot/internal/database"
	"me-bot/internal/handler"
	"me-bot/internal/repository"
	"me-bot/internal/service"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg)

	bot, err := linebot.New(cfg.LineChannelSecret, cfg.LineAccessToken)
	if err != nil {
		log.Fatal("LINE bot init error:", err)
	}

	userRepo := repository.NewUserRepository(db)
	attRepo := repository.NewAttendanceRepository(db)
	checkinSvc := service.NewCheckinService(bot, userRepo, attRepo, cfg)
	webhookHandler := handler.NewWebhookHandler(bot, checkinSvc)

	http.HandleFunc("/webhook", webhookHandler.Handle)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"ME Bot"}`))
	})

	log.Printf("🚀 ME Bot starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
