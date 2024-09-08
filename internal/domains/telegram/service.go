package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegrammbot.core/internal/config"
	"telegrammbot.core/internal/entities/sheet"
)

var (
	ErrInvalidValuesCount = errors.New("invalid values count")
)

type (
	ISheetService interface {
		HandleRequest(ctx context.Context, rawRequest string, reqType sheet.ReqType) (result string, err error)
	}
)

type Service struct {
	bot          *tgbotapi.BotAPI
	sheetService ISheetService
}

func NewService(sheetService ISheetService, botOpts config.BotOpts) (service *Service, err error) {
	bot, err := tgbotapi.NewBotAPI(botOpts.Token)
	if err != nil {
		return nil, fmt.Errorf("NewService: %w", err)
	}
	bot.Debug = true

	service = &Service{
		bot:          bot,
		sheetService: sheetService,
	}

	return service, nil
}

func (s *Service) Run(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			var msg tgbotapi.MessageConfig

			if strings.Contains(update.Message.Text, "/help") {
				result, err := s.sheetService.HandleRequest(ctx, update.Message.Text, sheet.Help)
				if err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, result)
				}
			} else if strings.Contains(update.Message.Text, "Добавить в") {
				if _, err := s.sheetService.HandleRequest(ctx, update.Message.Text, sheet.AddValueToCell); err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Успешно добавлено")
				}
			} else {
				result, err := s.sheetService.HandleRequest(ctx, update.Message.Text, sheet.GetValueFromCell)
				if err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Осталось потратить: %v рублей", result))
				}
			}
			s.bot.Send(msg)
		}
	}
}
