package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"io"
	"log"
	"os"
	"time"
)

type TgBot struct {
	bot   *bot.Bot
	store *ChatsStore
	ctx   context.Context
}

func (t *TgBot) Start() {
	t.bot.Start(t.ctx)
}

func (t *TgBot) RegisterHandlers() {
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, "/unexpected_wednesday", bot.MatchTypeExact, t.CommandHandlerUnexpectedWednesday)
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, "/iwantwednesday", bot.MatchTypeExact, t.CommandHandlerIWantWednesday)
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, "/wewantwednesday", bot.MatchTypeExact, t.CommandHandlerAddChat)
}

func (t *TgBot) CommandHandlerUnexpectedWednesday(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Да, внатуре, не среда"})
}

func (t *TgBot) CommandHandlerAddChat(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.Chat.Type == "private" {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Ты жадный до жаб, я не буду с тобой делиться"})
		return
	}

	if t.store.AddChat(update.Message.Chat.ID) {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Однажды и в нашем чате наступит среда"})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Этот чат уже получает жаб"})
	}
}

func (t *TgBot) CommandHandlerIWantWednesday(ctx context.Context, b *bot.Bot, update *models.Update) {
	message, _ := b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Ща, немного подожди"})
	defer b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID: message.Chat.ID, MessageID: message.ID,
	})

	file, err := os.CreateTemp(".", "*.jpeg")
	if err != nil {
		log.Println(err)
		return
	}
	defer func(file *os.File) {
		name := file.Name()
		file.Close()
		os.Remove(name)
	}(file)

	err = DownloadFrogPhoto(file)
	if err != nil {
		log.Println(err)
		return
	}
	file.Seek(0, 0)
	caption := "Браток, сегодня конечно не среда, но держи"
	if time.Now().Weekday() == time.Wednesday {
		caption = "Браток, сегодня среда, ты заслужил!"
	}
	t.SendImageReply(update.Message.Chat.ID, update.Message.ID, caption, file)
}

func (t *TgBot) SendImage(chatID int64, caption string, image io.Reader) {
	_, err := t.bot.SendPhoto(t.ctx, &bot.SendPhotoParams{
		ChatID:  chatID,
		Caption: caption,
		Photo: &models.InputFileUpload{
			Filename: "жаба.жпег",
			Data:     image,
		},
	})
	if err != nil {
		log.Println(err)
	}
}

func (t *TgBot) SendImageReply(chatID int64, replyTo int, caption string, image io.Reader) {
	_, err := t.bot.SendPhoto(t.ctx, &bot.SendPhotoParams{
		ChatID: chatID,
		ReplyParameters: &models.ReplyParameters{
			ChatID:    chatID,
			MessageID: replyTo,
		},
		Caption: caption,
		Photo: &models.InputFileUpload{
			Filename: "жаба.жпег",
			Data:     image,
		},
	})
	if err != nil {
		log.Println(err)
	}
}

func NewTgBot(token string, ctx context.Context) *TgBot {
	b, err := bot.New(token)
	if err != nil {
		panic(err)
	}

	return &TgBot{bot: b, ctx: ctx}
}
