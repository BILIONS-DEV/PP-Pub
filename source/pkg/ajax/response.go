package ajax

const (
	ERROR   = "error"
	SUCCESS = "success"
	WARNING = "warning"
)

type Responses struct {
	Status     string      `json:"status,omitempty"`
	Message    string      `json:"message,omitempty"`
	Errors     []Error     `json:"errors,omitempty"`
	DataObject interface{} `json:"data_object,omitempty"`
}

type Error struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}
