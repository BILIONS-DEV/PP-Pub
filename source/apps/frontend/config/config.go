package config

var (
	TitlePrefix = "Valueimpression"
)

// TitleWithPrefix render title with prefix default
//
// param: title
// return:
func TitleWithPrefix(title string) (titleWithPrefix string) {
	titleWithPrefix = title + " - " + TitlePrefix
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
