package aerospike

func (SetBlockedPage) SetName() string {
	return "blocked_page"
}

type SetBlockedPage struct {
	IsBlocked bool `as:"is_blocked"`
}
