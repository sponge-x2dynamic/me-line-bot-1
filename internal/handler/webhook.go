package handler

import (
	"log"
	"me-bot/internal/service"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type WebhookHandler struct {
	bot *linebot.Client
	svc *service.CheckinService
}

func NewWebhookHandler(bot *linebot.Client, svc *service.CheckinService) *WebhookHandler {
	return &WebhookHandler{bot: bot, svc: svc}
}

func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	events, err := h.bot.ParseRequest(r)
	if err != nil {
		log.Println("ParseRequest error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, event := range events {
		switch event.Type {

		case linebot.EventTypeMessage:
			if event.Source.Type != linebot.EventSourceTypeUser {
				continue
			}
			switch msg := event.Message.(type) {
			case *linebot.TextMessage:
				h.svc.HandleText(event, msg.Text)
			case *linebot.LocationMessage:
				h.svc.HandleLocation(event, msg.Latitude, msg.Longitude)
			}

		case linebot.EventTypeJoin:
			log.Printf("Bot joined group: %s", event.Source.GroupID)
		}
	}

	w.WriteHeader(http.StatusOK)
}
