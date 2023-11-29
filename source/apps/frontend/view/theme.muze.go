package view

type Muze struct {
	Logo   string
	Card   string
	CardBg string
}

type ThemeSetting struct {
	DevVersion  Muze
	ReleVersion Muze
}

var Setting ThemeSetting

func init() {
	MakeThemeSetting()
}

func MakeThemeSetting() {
	Setting.DevVersion = Muze{
		Card:   "card-primary",
		CardBg: "bg-primary-gradient",
	}
	Setting.ReleVersion = Muze{
		Card:   "card-dark",
		CardBg: "bg-info-gradient",
	}
}
