package lang

type HistoryError struct {
	List  TYPETranslation `json:"list"`
}

type History struct {
	Title TYPETranslation `json:"title"`
	Run   TYPETranslation `json:"run"`
}