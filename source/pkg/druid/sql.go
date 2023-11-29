package druid

const (
	// For ResultFormat
	FormatObject      = "object"
	FormatArray       = "array"
	FormatObjectLines = "objectLines"
	FormatArrayLines  = "arrayLines"
	FormatCsv         = "csv"

	// For Type of Parameter
	ParameterTypeCHAR     = "CHAR"
	ParameterTypeVARCHAR  = "VARCHAR"
	ParameterTypeDECIMAL  = "DECIMAL"
	ParameterTypeFLOAT    = "FLOAT"
	ParameterTypeREAL     = "REAL"
	ParameterTypeDOUBLE   = "DOUBLE"
	ParameterTypeBOOLEAN  = "BOOLEAN"
	ParameterTypeTINYINT  = "TINYINT"
	ParameterTypeSMALLINT = "SMALLINT"
	ParameterTypeINTEGER  = "INTEGER"
	ParameterTypeBIGINT   = "BIGINT"
	ParameterTIMESTAMP    = "TIMESTAMP"
	ParameterDATE         = "TIMESTAMP"
)

type Sql struct {
	Query        string                 `json:"query"`
	resultFormat string                 `json:"resultFormat"`
	Header       bool                   `json:"header"`
	Context      map[string]interface{} `json:"context"`
	Parameters   []Parameter            `json:"parameters"`
}

type Parameter struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
