package lang

type Channels struct {
	Title          TYPETranslation `json:"title"`
	Add            TYPETranslation `json:"add"`
	Edit           TYPETranslation `json:"edit"`
	SearchChannels TYPETranslation `json:"search_channels"`
	Run            TYPETranslation `json:"run"`
	Top            TYPETranslation `json:"top"`
	Main           TYPETranslation `json:"main"`
	Name           TYPETranslation `json:"name"`
	Description    TYPETranslation `json:"description"`
	Category       TYPETranslation `json:"category"`
	Keyword        TYPETranslation `json:"keyword"`
	Language       TYPETranslation `json:"language"`
	Button         TYPETranslation `json:"button"`
}

type ChannelsError struct {
	Add      TYPETranslation `json:"add"`
	Edit     TYPETranslation `json:"edit"`
	Delete   TYPETranslation `json:"delete"`
	List     TYPETranslation `json:"list"`
	NotFound TYPETranslation `json:"not_found"`
}
