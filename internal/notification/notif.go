package notification

import (
	"context"
	"fmt"
	"gabutlabs/devopin-cli/internal/config"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Notification struct {
	ctx context.Context
	cfg *config.Config
}

func NewNotification(ctx context.Context, cfg *config.Config) *Notification {
	return &Notification{
		ctx: ctx,
		cfg: cfg,
	}
}

func (n *Notification) FormatResourceAlertMessage(hostName string, resource string, value float64, threshold int) string {
	message := fmt.Sprintf(
		"<b>%s Alert</b>\n\n"+
			"Server : <code>%s</code>\n"+
			"Usage  : <code>%.2f%%</code>\n"+
			"Limit  : <code>%d%%</code>\n"+
			"Time   : <code>%s</code>\n",
		resource,
		strings.ToUpper(hostName),
		value,
		threshold,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	return message
}

func (n *Notification) SendTelegramAlert(message string) {
	b, err := bot.New(n.cfg.Notify.Telegram.BotToken)
	if err != nil {
		panic(err)
	}

	// Implementasi pengiriman notifikasi Telegram
	_, err = b.SendMessage(n.ctx, &bot.SendMessageParams{
		ChatID:    n.cfg.Notify.Telegram.ChatID,
		Text:      message,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		panic(err)
	}
}
