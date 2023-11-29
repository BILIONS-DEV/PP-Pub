package mysql

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type RuleMetrics []RuleMetric
type RuleMetric struct {
	Name      string  `json:"name"`
	Compare   string  `json:"compare"`
	Threshold float64 `json:"threshsold"`
}

// Value return json value, implement driver.Valuer interface
func (m RuleMetrics) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	ba, err := json.Marshal(m)

	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *RuleMetrics) Scan(val interface{}) error {
	if val == nil {
		*m = RuleMetrics{}
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal value:", val))
	}
	err := json.Unmarshal(ba, m)
	return err
}

func (m RuleMetrics) String() string {
	a, _ := json.Marshal(m)
	return string(a)
}

// GormDataType gorm common data type
func (RuleMetrics) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (RuleMetrics) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

func (m RuleMetrics) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if len(m) == 0 {
		return gorm.Expr("NULL")
	}

	data, _ := json.Marshal(m)

	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}

type ObservedDimensions []ObservedDimension
type ObservedDimension struct {
	Key       string   `json:"key"`
	Value     []string `json:"value,omitempty"`
	FilterAll bool     `json:"filterAll,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m ObservedDimensions) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	ba, err := json.Marshal(m)

	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *ObservedDimensions) Scan(val interface{}) error {
	if val == nil {
		*m = ObservedDimensions{}
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal value:", val))
	}
	err := json.Unmarshal(ba, m)
	return err
}

func (m ObservedDimensions) String() string {
	a, _ := json.Marshal(m)
	return string(a)
}

// GormDataType gorm common data type
func (ObservedDimensions) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (ObservedDimensions) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

func (m ObservedDimensions) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if len(m) == 0 {
		return gorm.Expr("NULL")
	}

	data, _ := json.Marshal(m)

	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}

type NotifyUsers []int64

// Value return json value, implement driver.Valuer interface
func (m NotifyUsers) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	ba, err := json.Marshal(m)

	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *NotifyUsers) Scan(val interface{}) error {
	if val == nil {
		*m = NotifyUsers{}
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal value:", val))
	}
	err := json.Unmarshal(ba, m)
	return err
}

func (m NotifyUsers) String() string {
	a, _ := json.Marshal(m)
	return string(a)
}

// GormDataType gorm common data type
func (NotifyUsers) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (NotifyUsers) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

func (m NotifyUsers) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if len(m) == 0 {
		return gorm.Expr("NULL")
	}

	data, _ := json.Marshal(m)

	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}
