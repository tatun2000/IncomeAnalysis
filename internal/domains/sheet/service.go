package sheet

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"telegrammbot.core/internal/config"
	"telegrammbot.core/internal/entities/sheet"
	"telegrammbot.core/internal/errs"
)

type (
	IOauthService interface {
		GetClient() *http.Client
	}

	Service struct {
		oauthService IOauthService

		spreedSheetID string
	}
)

func NewService(oauthService IOauthService, sheetOpts config.SheetOpts) *Service {
	return &Service{
		oauthService:  oauthService,
		spreedSheetID: sheetOpts.SpreadsheetId,
	}
}

func (s *Service) HandleRequest(ctx context.Context, rawRequest string, reqType sheet.ReqType) (result string, err error) {
	switch reqType {
	case sheet.AddValueToCell:
		reqItems := strings.Split(rawRequest, " ")
		if len(reqItems) != 4 {
			return result, fmt.Errorf("HandleRequest: %w", errs.ErrInvalidRequestMessageFormat)
		}

		value, err := strconv.ParseFloat(reqItems[1], 32)
		if err != nil {
			return result, fmt.Errorf("HandleRequest: %w", err)
		}

		categoryNum, err := strconv.Atoi(reqItems[3])
		if err != nil {
			return result, fmt.Errorf("HandleRequest: %w", err)
		}

		month, err := sheet.GetActualMonthSheet()
		if err != nil {
			return result, fmt.Errorf("HandleRequest: %w", err)
		}
		day, err := sheet.GetActualDayCell()
		if err != nil {
			return result, fmt.Errorf("HandleRequest: %w", err)
		}
		category, err := sheet.ConvertCategoryTypeToCell(categoryNum)
		if err != nil {
			return result, fmt.Errorf("HandleRequest: %w", err)
		}

		tablePath := fmt.Sprintf("%s!%s%s", month, day, category)
		log.Printf("tablePath = %s", tablePath)

		client := s.oauthService.GetClient()
		srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			return result, fmt.Errorf("HandleRequest: %w", err)
		}

		valueRange := &sheets.ValueRange{
			Values: [][]interface{}{
				{value},
			},
		}
		if _, err = srv.Spreadsheets.Values.Update(s.spreedSheetID, tablePath, valueRange).ValueInputOption("RAW").Do(); err != nil {
			return result, fmt.Errorf("HandleRequest: %w", err)
		}
	case sheet.GetValueFromCell:
		client := s.oauthService.GetClient()
		srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			return result, fmt.Errorf("readTotalAmountToSpend: %w", err)
		}
		resp, err := srv.Spreadsheets.Values.Get(s.spreedSheetID, "Август!AH38:AJ39").Do()
		if err != nil {
			return result, fmt.Errorf("readTotalAmountToSpend: %w", err)
		}

		if len(resp.Values) == 0 {
			return "", nil
		} else {
			for _, row := range resp.Values {
				result = fmt.Sprintf("%s", row[0])
			}
		}
	default:
		return result, fmt.Errorf("HandleRequest: %w", errs.ErrUnknownRequestType)
	}
	return
}
