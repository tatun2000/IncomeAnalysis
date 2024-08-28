package oauth

type (
	getRefreshTokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}

	tokenStruct struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
	}
)
