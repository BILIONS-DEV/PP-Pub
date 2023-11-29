package druid

type client struct {
	API string
}

func NewClient(api string, args ...string) *client {
	return &client{API: api}
}

func (t *client) Execute(query string, output interface{}) (err error) {
	return
}
