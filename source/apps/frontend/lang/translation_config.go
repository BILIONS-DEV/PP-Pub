package lang

type Config struct {
	Title              TYPETranslation `json:"title"`
	PrebidTimeOut      TYPETranslation `json:"prebid_time_out"`
	AdRefreshTime      TYPETranslation `json:"ad_refresh_time"`
	SelectCurrency     TYPETranslation `json:"select_currency"`
	SelectCurrencyDesc TYPETranslation `json:"select_currency_description"`
	Save               TYPETranslation `json:"save"`
}

type ConfigError struct {
	Save TYPETranslation `json:"save"`
}
