package druid

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"reflect"
	"time"
)

var Client *Druid

type Druid struct {
	url       string
	collector *colly.Collector
	req       *colly.Request
	res       *colly.Response
	sql       Sql
}

func NewClient(url string) *Druid {
	c := colly.NewCollector()
	return &Druid{
		url:       url,
		collector: c,
	}
}

func (this *Druid) Find(a interface{}) (err error) {
	rv := reflect.ValueOf(a)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	} else {
		err = errors.New("require is a pointer")
		return
	}
	res, err := this.requestSql()
	if err != nil {
		return
	}
	err = this.changerv(rv, res, -1)
	return
}

func (this *Druid) changerv(rv reflect.Value, res []map[string]interface{}, i int) (err error) {
	if rv.Kind() == reflect.Struct {
		this.changeStruct(rv, res, i)
	}
	if rv.Kind() == reflect.Slice {
		err = this.changeSlice(rv, res)
	}
	return
}

// assumes rv is a slice
func (this *Druid) changeSlice(rv reflect.Value, res []map[string]interface{}) (err error) {
	ln := rv.Len()
	if ln == 0 && rv.CanAddr() {
		var elem reflect.Value

		typ := rv.Type().Elem()
		if typ.Kind() == reflect.Ptr {
			elem = reflect.New(typ.Elem())
		}
		if typ.Kind() == reflect.Struct {
			elem = reflect.New(typ).Elem()
		}
		for i := 0; i < len(res); i++ {
			rv.Set(reflect.Append(rv, elem))
		}
	}
	//fmt.Println(len(res))
	ln = rv.Len()
	//fmt.Println(ln)
	for i := 0; i < ln; i++ {
		err = this.changerv(rv.Index(i), res, i)
		if err != nil {
			return
		}
	}
	return
}

// assumes rv is a struct
func (this *Druid) changeStruct(rv reflect.Value, res []map[string]interface{}, indexRes int) {
	if !rv.CanAddr() {
		return
	}
	for i := 0; i < rv.NumField(); i++ {
		//fmt.Printf("%+v \n", res)
		//fmt.Println(i)
		field := rv.Field(i)
		value := res[indexRes][rv.Type().Field(i).Tag.Get("json")]
		if value == nil {
			continue
		}
		//fmt.Println(rv.Type().Field(i).Tag.Get("json"))
		//fmt.Println(field.Type().Kind())
		//fmt.Println(reflect.ValueOf(value).Type())
		//fmt.Println("value",value)
		switch _ := field.Interface().(type) {
		case time.Time:
			layout := "2006-01-02T15:04:05.000Z"
			t, err := time.Parse(layout, value.(string))
			if err != nil {
				fmt.Println(err)
				continue
			}
			field.Set(reflect.ValueOf(t))
			break
		case float64:
			field.SetFloat(value.(float64))
			break
		case int:
			switch valueRes := value.(type) {
			case float64:
				field.SetInt(int64(valueRes))
			default:
				field.SetInt(valueRes.(int64))
			}
			break
		default:
			field.Set(reflect.ValueOf(value))
		}
		//switch field.Kind() {
		//case reflect.String:
		//	str := rv.Type().Field(i).Name
		//	fmt.Println(str)
		//	//tag := ParseTagSetting(rv.Type().Field(i).Tag.Get("druid"),":")
		//	fmt.Println(rv.Type().Field(i).Tag.Get("json"))
		//	field.SetString("fred")
		//case reflect.Int:
		//	field.SetInt(54)
		//case reflect.Float64:
		//	field.SetFloat(54)
		//case reflect.TypeOf(time.Time{}).Kind():
		//	return v.String()
		//default:
		//	fmt.Println("unknown field")
		//}
	}
}

func (this *Druid) requestSql() (res []map[string]interface{}, err error) {
	defer func() { this.sql = Sql{} }()
	if this.sql.Query == "" {
		err = errors.New("query is require")
		return
	}
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/json")
	})
	// authenticate
	// attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &res)
		if err != nil {
			return
		}
	})

	this.sql.resultFormat = FormatObject
	byteSql, _ := json.Marshal(this.sql)
	if err != nil {
		return
	}

	fmt.Println(string(byteSql))
	//err = c.PostRaw("http://66.206.12.106:8888/druid/v2/sql", []byte(qry))
	err = c.PostRaw(this.url+URIApiSql, byteSql)
	if err != nil {
		log.Fatal("post: ", err)
	}
	//fmt.Println(res)
	//var i interface{}
	//i = &res
	//rv := reflect.ValueOf(i)
	////fmt.Println(string(rv.Bytes()))
	//for rv.Kind() == reflect.Ptr {
	//	if rv.IsNil() && rv.CanAddr() {
	//		rv.Set(reflect.New(rv.Type().Elem()))
	//	}
	//	rv = rv.Elem()
	//}
	return
}

func (this *Druid) Sql(sql Sql) *Druid {
	this.sql = sql
	return this
}
