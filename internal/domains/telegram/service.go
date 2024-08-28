package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegrammbot.core/internal/constants"
)

var (
	ErrInvalidValuesCount = errors.New("invalid values count")
)

type (
	IOauthService interface {
		GetClient() *http.Client
	}

	ISheetService interface {
		FetchCellsValue(spreadsheetID, rangeCells string) (result []string, err error)
	}
)

type Service struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	bot          *tgbotapi.BotAPI
	oauthService IOauthService
	sheetService ISheetService
}

func NewService(oauthService IOauthService, sheetService ISheetService, token string) (service *Service, err error) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("NewService: %w", err)
	}
	bot.Debug = true

	service = &Service{
		ctx:        ctx,
		cancelFunc: cancelFunc,

		bot:          bot,
		oauthService: oauthService,
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
			switch update.Message.Text {
			case "Сколько осталось потратить?":
				result, err := s.readTotalAmountToSpend(ctx)
				if err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
					break
				}
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, result)
			default:
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Не то")
			}
			s.bot.Send(msg)
		}
	}
}

func (s *Service) readTotalAmountToSpend(ctx context.Context) (result string, err error) {
	//client := s.oauthService.GetClient()

	values, err := s.sheetService.FetchCellsValue(constants.SpreadsheetID, "Август!AH38:AJ39")
	if err != nil {
		return result, fmt.Errorf("readTotalAmountToSpend: %w", err)
	}

	if len(values) != 1 {
		return result, fmt.Errorf("readTotalAmountToSpend: %w", &ErrInvalidValuesCount)
	}

	return fmt.Sprintf("Осталось потратить: %v рублей", result), nil
}
