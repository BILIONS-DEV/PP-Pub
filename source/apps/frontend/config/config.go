package config

var (
	TitlePrefix = "Your Brand"
)

// TitleWithPrefix render title with prefix default
//
// param: title
// return:
func TitleWithPrefix(title string) (titleWithPrefix string) {
	if TitlePrefix != "" {
		titleWithPrefix = title + " - " + TitlePrefix
	} else {
		titleWithPrefix = title
	}
	return
}

const (
	STATUSError   = "error"
	STATUSSuccess = "success"
	STATUSWarning = "warning"
)

const (
	JwtSecret      = "my_secret"
	JwtContextKey  = "jwtToken"
	JwtCookieName  = "jateox"
	JwtTokenLookup = "cookie:" + JwtCookieName
)
