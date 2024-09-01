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
	ctx        context.Context
	cancelFunc context.CancelFunc

	bot          *tgbotapi.BotAPI
	sheetService ISheetService
}

func NewService(sheetService ISheetService, botOpts config.BotOpts) (service *Service, err error) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	bot, err := tgbotapi.NewBotAPI(botOpts.Token)
	if err != nil {
		return nil, fmt.Errorf("NewService: %w", err)
	}
	bot.Debug = true

	service = &Service{
		ctx:        ctx,
		cancelFunc: cancelFunc,

		bot:          bot,
		sheetService: sheetService,
	}

	go service.handleRequests(service.ctx)

	return service, nil
}

func (s *Service) Close() {
	s.ctx.Done()
}

func (s *Service) handleRequests(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			var msg tgbotapi.MessageConfig

			if strings.Contains(update.Message.Text, "Добавить в") {
				if _, err := s.sheetService.HandleRequest(ctx, update.Message.Text, sheet.AddValueToCell); err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
					break
				}
			} else {
				result, err := s.sheetService.HandleRequest(ctx, update.Message.Text, sheet.GetValueFromCell)
				if err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
					break
				}
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Осталось потратить: %v рублей", result))
			}
			s.bot.Send(msg)
		}
	}
}
