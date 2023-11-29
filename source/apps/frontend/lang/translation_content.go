package lang

type Content struct {
	Title         TYPETranslation `json:"title"`
	Add           TYPETranslation `json:"add"`
	Edit          TYPETranslation `json:"edit"`
	SearchContent TYPETranslation `json:"search_content"`
	Run           TYPETranslation `json:"run"`
	Top           TYPETranslation `json:"top"`
	Main          TYPETranslation `json:"main"`
	Name          TYPETranslation `json:"name"`
	Description   TYPETranslation `json:"description"`
	Video         TYPETranslation `json:"video"`
	PreviewVideo  TYPETranslation `json:"preview_video"`
	Thumb         TYPETranslation `json:"thumb"`
	PreviewThumb  TYPETranslation `json:"preview_thumb"`
	SelectThumb   TYPETranslation `json:"select_thumb"`
	Channels      TYPETranslation `json:"channels"`
	Category      TYPETranslation `json:"category"`
	Tags          TYPETranslation `json:"tags"`
	Keyword       TYPETranslation `json:"keyword"`
	VideoType     TYPETranslation `json:"video_type"`
	Button        TYPETranslation `json:"button"`
}

type ContentError struct {
	Add      TYPETranslation `json:"add"`
	Edit     TYPETranslation `json:"edit"`
	Delete   TYPETranslation `json:"delete"`
	List     TYPETranslation `json:"list"`
	NotFound TYPETranslation `json:"not_found"`
}
