package lang

type Playlist struct {
	Title          TYPETranslation `json:"title"`
	Add            TYPETranslation `json:"add"`
	Edit           TYPETranslation `json:"edit"`
	SearchPlaylist TYPETranslation `json:"search_playlist"`
	Run            TYPETranslation `json:"run"`
	Top            TYPETranslation `json:"top"`
	Main           TYPETranslation `json:"main"`
	Name           TYPETranslation `json:"name"`
	Description    TYPETranslation `json:"description"`
	User           TYPETranslation `json:"user"`
	Channels       TYPETranslation `json:"channels"`
	Language       TYPETranslation `json:"language"`
	Category       TYPETranslation `json:"category"`
	Keywords       TYPETranslation `json:"keywords"`
	Content        TYPETranslation `json:"content"`
	Videos         TYPETranslation `json:"videos"`
	SelectContent  TYPETranslation `json:"select_content"`
	Button         TYPETranslation `json:"button"`
	PlaylistConfig TYPETranslation `json:"playlist_config"`
	Clear          TYPETranslation `json:"clear"`
}

type PlaylistError struct {
	Add               TYPETranslation `json:"add"`
	Edit              TYPETranslation `json:"edit"`
	Delete            TYPETranslation `json:"delete"`
	List              TYPETranslation `json:"list"`
	NotFound          TYPETranslation `json:"not_found"`
	RlContentPlaylist TYPETranslation `json:"rl_content_playlist"`
	PlaylistConfig    TYPETranslation `json:"playlist_config"`
}
