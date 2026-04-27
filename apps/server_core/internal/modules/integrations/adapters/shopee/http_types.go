package shopee

type shopeeAPIError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type shopeeTokenResponse struct {
	RequestID string `json:"request_id"`
	Error     string `json:"error"`
	Message   string `json:"message"`
	Response  struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpireIn     int    `json:"expire_in"`
		ShopID       string `json:"shop_id"`
	} `json:"response"`
}
