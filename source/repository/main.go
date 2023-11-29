package main

import (
	"fmt"
	"source/repository/payment/invoice"
)

// TungDT viết phần này để chạy invoice hôm mồng 3 tết ( có facetime)
func main() {
	err := new(invoice.PaymentRequest).ScanPaymentPublisher()
	fmt.Printf("\n err: %+v \n", err)
}
