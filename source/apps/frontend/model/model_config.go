package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
	"source/apps/frontend/payload"
	"source/core/technology/mysql"
)

type Config struct{}

type ConfigRecord struct {
	mysql.TableConfig
}

func (ConfigRecord) TableName() string {
	return mysql.Tables.Config
}

func (t *Config) GetByUserId(id int64) (record ConfigRecord) {
	mysql.Client.Where(ConfigRecord{mysql.TableConfig{UserId: id}}).Find(&record)
	return
}

func (t *Config) ReadCurrencyPrebid() (resp payload.ResponseCurrency, err error, statusCode int) {
	url := "https://cdn.jsdelivr.net/gh/prebid/currency-file@1/latest.json"
	//makeRequest(c)
	var c = colly.NewCollector(
		colly.UserAgent("get currency by APD"),
		colly.AllowURLRevisit(),
	)

	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &resp)
	})
	// Set error handler
	c.OnError(func(r *colly.Response, errHandle error) {
		statusCode = r.StatusCode
		if errHandle != nil {
			err = errHandle
		} else {
			if r.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("status code error: %d", r.StatusCode))
			}
		}
	})
	errVisit := c.Visit(url)
	if err == nil {
		err = errVisit
	}
	c.Wait()
	return
}

func (config *ConfigRecord) IsSetCurrency() bool {
	if config.Currency == "" {
		return false
	}
	return true
}

func (t *Config) GetSymbolCurrencyByUserId(id int64) (currency, symbol string) {
	record := ConfigRecord{}
	mysql.Client.Where(ConfigRecord{mysql.TableConfig{UserId: id}}).Find(&record)
	currencyInfo := new(Currency).GetByCode(record.Currency)
	currency = record.Currency
	symbol = currencyInfo.Symbol
	return
}
