package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const refreshToken = "1//0cDHXAVH1wIE4CgYIARAAGAwSNwF-L9Ircyl17EGWJzg_E5ApKOCUA0X2aEGSm2kHQcoHhr7BUUM6TIDSw9_D6kgiouHkqWHp7pI"
const spreadsheetId = "1T28n-GhmvDXeTICR3X7gzvPtU-bozmv_Ys_Uid5WK80"

var Cfg *oauth2.Config

type Credentials struct {
	Type                string `json:"type"`
	ProjectID           string `json:"project_id"`
	PrivateKeyID        string `json:"private_key_id"`
	PrivateKey          string `json:"private_key"`
	ClientEmail         string `json:"client_email"`
	ClientID            string `json:"client_id"`
	AuthURI             string `json:"auth_uri"`
	TokenURI            string `json:"token_uri"`
	AuthProviderCertURL string `json:"auth_provider_x509_cert_url"`
	ClientCertURL       string `json:"client_x509_cert_url"`
}

func main() {
	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt, syscall.SIGTERM)
	defer cancelFunc()

	app, gracefulShutdown, err := InjectAppGod(ctx, ".")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		app.telegramService.Run(ctx)
	}()

	<-ctx.Done()
	gracefulShutdown()
}

func oldMain() {
	ctx := context.Background()
	b, err := os.ReadFile("../credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	if err = updateAccessToken(ctx); err != nil {
		log.Fatal(err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	config.RedirectURL = "http://localhost:8080/callback"
	Cfg = config

	// http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
	// 	code := r.URL.Query().Get("code")
	// 	if code == "" {
	// 		http.Error(w, "Код не найден", http.StatusBadRequest)
	// 		return
	// 	}
	// 	fmt.Fprintf(w, "code: %v\n", code)
	// 	// Обмен кода на токен
	// 	token, err := config.Exchange(context.Background(), code)
	// 	if err != nil {
	// 		http.Error(w, "Не удалось получить токен", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	fmt.Fprintf(w, "Токен доступа: %v\n", token.AccessToken)
	// 	fmt.Fprintf(w, "Refresh токен: %v\n", token.RefreshToken)
	// 	fmt.Fprintf(w, "тип: %v\n", token.TokenType)

	// 	fmt.Fprintf(w, "exp: %v", token.Expiry)
	// })

	// go func() {
	// 	log.Printf("Запуск сервера на :8080")
	// 	log.Fatal(http.ListenAndServe(":8080", nil))
	// }()

	botRun(ctx)
}

func readTotalAmountToSpend(ctx context.Context) (result string, err error) {
	if Cfg == nil {
		return result, errors.New("Cfg is empty")
	}
	client := getClient(Cfg)

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
			result = fmt.Sprintf("%s", row[0])
		}
	}

	return fmt.Sprintf("Осталось потратить: %v рублей", result), nil
}

func botRun(ctx context.Context) {
	// / telegram bot
	bot, err := tgbotapi.NewBotAPI("5657441460:AAHu_VJ3jBt9Nv2uXvMkuabkXkkvaD70oA8")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			var msg tgbotapi.MessageConfig
			switch update.Message.Text {
			case "Сколько осталось потратить?":
				result, err := readTotalAmountToSpend(ctx)
				if err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
					break
				}
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, result)
			default:
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Не то")
			}
			bot.Send(msg)
		}
	}
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "../token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func updateAccessToken(ctx context.Context) (err error) {

	client := resty.New()

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}

	// Выполняем POST-запрос с параметрами в теле, используя FormData
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"client_id":     "228230793527-vhj698i1n7m6i6nietr666235sfs6aik.apps.googleusercontent.com",
			"client_secret": "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",
			"refresh_token": refreshToken,
			"grant_type":    "refresh_token",
		}).
		Post("https://oauth2.googleapis.com/token")

	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}

	// Обрабатываем ответ
	fmt.Printf("Статус-код: %d\n", resp.StatusCode())
	fmt.Printf("Ответ: %s\n", resp.Body())

	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("updateAccessToken: %w", err)
	}

	tokenStruct := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
	}{
		AccessToken:  result.AccessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}

	// Создаем файл
	file, err := os.Create("../token.json")
	if err != nil {
		return fmt.Errorf("updateAccessToken: %w", err)
	}
	defer file.Close()

	// Кодируем структуру в JSON и записываем в файл
	encoder := json.NewEncoder(file)
	err = encoder.Encode(tokenStruct)
	if err != nil {
		return fmt.Errorf("updateAccessToken: %w", err)
	}

	return nil
}
