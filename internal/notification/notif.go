package notification

import (
	"context"
	"fmt"
	"gabutlabs/devopin-cli/internal/config"
	"log"
	"net/smtp"
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

func (n *Notification) FormatMonitorWorkerAlertMessage(hostName string, workers []string) string {
	var workerList strings.Builder

	for _, w := range workers {
		fmt.Fprintf(&workerList, "• <code>%s</code>\n", w)
	}

	message := fmt.Sprintf(
		"<b>Worker Alert</b>\n\n"+
			"<b>Server</b> : <code>%s</code>\n"+
			"<b>Problematic Workers</b> : <code>%d</code>\n\n"+
			"%s\n"+
			"<b>Time</b> : <code>%s</code>\n",

		strings.ToUpper(hostName),
		len(workers),
		workerList.String(),
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return message

}

// SendNotif sends notifications to all enabled channels
func (n *Notification) SendNotif(message string, subject string) {
	if n.cfg.Notify.Channels.Telegram {
		n.sendTelegramAlert(message)
	}
	if n.cfg.Notify.Channels.Email {
		n.sendEmailAlert(message, subject)
	}
}

func (n *Notification) sendTelegramAlert(message string) {
	b, err := bot.New(n.cfg.Notify.Telegram.BotToken)
	if err != nil {
		log.Printf("[Telegram] failed to init bot: %v", err)
		return
	}

	_, err = b.SendMessage(n.ctx, &bot.SendMessageParams{
		ChatID:    n.cfg.Notify.Telegram.ChatID,
		Text:      message,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		log.Printf("[Telegram] failed to send message: %v", err)
		return
	}

	log.Printf("[Telegram] notification sent successfully")
}

func (n *Notification) sendEmailAlert(message string, subject string) {
	if len(n.cfg.Notify.Email.ToEmails) == 0 {
		log.Printf("[Email] no recipient emails configured")
		return
	}

	// Build email message
	msg := fmt.Sprintf(
		"To: %s\r\n"+
			"From: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		strings.Join(n.cfg.Notify.Email.ToEmails, ", "),
		n.cfg.Notify.Email.FromEmail,
		subject,
		message,
	)

	auth := smtp.PlainAuth(
		"",
		n.cfg.Notify.Email.SMTPUser,
		n.cfg.Notify.Email.SMTPPassword,
		n.cfg.Notify.Email.SMTPHost,
	)
	addr := fmt.Sprintf("%s:%d", n.cfg.Notify.Email.SMTPHost, n.cfg.Notify.Email.SMTPPort)

	// Send to all recipients
	err := smtp.SendMail(
		addr,
		auth,
		n.cfg.Notify.Email.FromEmail,
		n.cfg.Notify.Email.ToEmails,
		[]byte(msg),
	)
	log.Println(err)
	if err != nil {
		log.Printf("[Email] failed to send email: %v", err)
		return
	}

	log.Printf("[Email] notification sent successfully to %d recipient(s)", len(n.cfg.Notify.Email.ToEmails))
}
