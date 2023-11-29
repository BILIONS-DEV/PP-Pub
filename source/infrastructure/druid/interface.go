package druid

type Druid interface {
	Execute(query string, output interface{}) (err error)
}
