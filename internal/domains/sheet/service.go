package sheet

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const spreadsheetId = "1T28n-GhmvDXeTICR3X7gzvPtU-bozmv_Ys_Uid5WK80"

type (
	IOauthService interface {
		GetClient() *http.Client
	}

	Service struct {
		oauthService IOauthService
	}
)

func NewService(oauthService IOauthService) *Service {
	return &Service{
		oauthService: oauthService,
	}
}

func (s *Service) FetchCellsValue(ctx context.Context, spreadsheetID, rangeCells string) (result []string, err error) {
	client := s.oauthService.GetClient()
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return result, fmt.Errorf("readTotalAmountToSpend: %w", err)
	}
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, "Август!AH38:AJ39").Do()
	if err != nil {
		return result, fmt.Errorf("readTotalAmountToSpend: %w", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for _, row := range resp.Values {
			result = append(result, fmt.Sprintf("%s", row[0]))
		}
	}

	return result, nil
}
