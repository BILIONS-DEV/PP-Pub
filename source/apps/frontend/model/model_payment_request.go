package model

import (
	"source/core/technology/mysql"
	"strconv"
	"strings"
)

type PaymentRequest struct{}

type PaymentRequestRecord struct {
	mysql.TablePaymentRequest
}

func (PaymentRequestRecord) TableName() string {
	return mysql.Tables.PaymentRequest
}

func (t *PaymentRequest) GetByInvoice(invoice PaymentRecord) (records []PaymentRequestRecord, err error) {
	ids := StringToArrayInt(invoice.RequestId, ",")
	err = mysql.Client.Where("id in ?", ids).Find(&records).Error
	return
}

func StringToArrayInt(string, text string) []int {
	strs := strings.Split(string, text)
	array := make([]int, len(strs))
	for i := range array {
		array[i], _ = strconv.Atoi(strs[i])
	}

	return array
}
