package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"telegrammbot.core/internal/config"
	"telegrammbot.core/internal/constants"
)

type (
	Service struct {
		ctx        context.Context
		cancelFunc context.CancelFunc

		refreshToken      string
		workDir           string
		tokenFullFilename string
		httpClient        *resty.Client
		oauth2Cfg         *oauth2.Config
	}
)

func NewService(generalOpts config.GeneralOpts, oauthOpts config.OauthOpts) (service *Service, cleanup func(), err error) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	service = &Service{
		ctx:        ctx,
		cancelFunc: cancelFunc,

		workDir:           generalOpts.WorkDir,
		refreshToken:      oauthOpts.RefreshToken,
		tokenFullFilename: filepath.Join(generalOpts.WorkDir, constants.TokenFilename),
		httpClient:        resty.New(),
	}

	if err = service.setOauth2Config(service.ctx); err != nil {
		return service, nil, fmt.Errorf("NewService: %w", err)
	}

	go service.refreshAccessTokenInBackground(service.ctx)

	return service, func() {
		service.Close()
	}, nil
}

func (s *Service) Close() {
	s.ctx.Done()
}

func (s *Service) setOauth2Config(ctx context.Context) (err error) {
	data, err := os.ReadFile(path.Join(s.workDir, constants.CredFilename))
	if err != nil {
		return fmt.Errorf("setOauth2Config: %w", err)
	}

	config, err := google.ConfigFromJSON(data, constants.SpreadsheetsScopeURL)
	if err != nil {
		return fmt.Errorf("setOauth2Config: %w", err)
	}
	config.RedirectURL = "http://localhost:8080/callback"
	s.oauth2Cfg = config
	return nil
}

func (s *Service) refreshAccessTokenInBackground(ctx context.Context) (err error) {
	if err = s.updateAccessToken(ctx); err != nil {
		return fmt.Errorf("refreshAccessTokenInBackground: %w", err)
	}
	ticker := time.NewTicker(time.Hour)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err = s.updateAccessToken(ctx); err != nil {
				return fmt.Errorf("refreshAccessTokenInBackground: %w", err)
			}
		}
	}
}

func (s *Service) updateAccessToken(ctx context.Context) (err error) {
	var result getRefreshTokenResp
	resp, err := s.httpClient.R().
		SetQueryParams(map[string]string{
			"client_id":     constants.ClientID,
			"client_secret": constants.ClientSecret,
			"refresh_token": s.refreshToken,
			"grant_type":    "refresh_token",
		}).
		Post(constants.FetchAccessTokenURL)
	if err != nil {
		return fmt.Errorf("updateAccessToken: %w", err)
	}

	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("updateAccessToken: %w", err)
	}

	token := tokenStruct{
		AccessToken:  result.AccessToken,
		RefreshToken: s.refreshToken,
		TokenType:    "Bearer",
	}

	file, err := os.Create(path.Join(s.workDir, constants.TokenFilename))
	if err != nil {
		return fmt.Errorf("updateAccessToken: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(token)
	if err != nil {
		return fmt.Errorf("updateAccessToken: %w", err)
	}

	return nil
}

func (s *Service) GetClient() (client *http.Client, err error) {
	file, err := os.Open(s.tokenFullFilename)
	if err != nil {
		return client, fmt.Errorf("GetClient: %w", err)
	}
	defer file.Close()

	token := &oauth2.Token{}
	if err = json.NewDecoder(file).Decode(token); err != nil {
		return client, fmt.Errorf("GetClient: %w", err)
	}

	return s.oauth2Cfg.Client(s.ctx, token), nil
}
