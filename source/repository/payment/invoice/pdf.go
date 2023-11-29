package invoice

import (
	"bytes"
	"html/template"
	"source/pkg/utility"
	"strconv"
	"time"

	// "github.com/andrewcharlton/wkhtmltopdf-go"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"source/core/technology/mysql"
)

type dataPDF struct {
	Invoice  PaymentInvoiceRecord
	Requests []PaymentRequestRecord
	User     UserRecord
	Billing  UserBillingRecord
	Address  string
	Today    string
}

// func GeneratePDF(data *domain.SomeModel) ([]byte, error) {
func GeneratePDF(invoice PaymentInvoiceRecord, paymentRequests []PaymentRequestRecord, user UserRecord, billing UserBillingRecord) (pathPDF string, err error) {
	var templ *template.Template

	var data = dataPDF{
		Invoice:  invoice,
		Requests: paymentRequests,
		User:     user,
		Billing:  billing,
		Address:  "",
		Today:    time.Now().Format("01/02/2006"),
	}
	data.Address = setAddress(user.Address, user.City, user.Country)
	// use Go's default HTML template generation tools to generate your HTML
	var body bytes.Buffer

	// initalize a wkhtmltopdf generator
	if !utility.IsWindow() {
		wkhtmltopdf.SetPath("wkhtmltopdf")
		if invoice.Status == mysql.StatusPaymentDone {
			if templ, err = template.ParseFiles("/home/assyrian/go/selfserve/source/repository/payment/invoice/pdf_invoice.gohtml"); err != nil {
				return
			}
		} else {
			if templ, err = template.ParseFiles("/home/assyrian/go/selfserve/source/repository/payment/invoice/pdf_statement.gohtml"); err != nil {
				return
			}
		}
	} else {
		wkhtmltopdf.SetPath("C:/Program Files/wkhtmltopdf/bin/wkhtmltopdf.exe")
		if invoice.Status == mysql.StatusPaymentDone {
			if templ, err = template.ParseFiles("D:/ProjectGo/PubPower/source/repository/payment/invoice/pdf_invoice.gohtml"); err != nil {
				return
			}
		} else {
			if templ, err = template.ParseFiles("D:/ProjectGo/PubPower/source/repository/payment/invoice/pdf_statement.gohtml"); err != nil {
				return
			}
		}
	}
	// apply the parsed HTML template data and keep the result in a Buffer
	if err = templ.Execute(&body, data); err != nil {
		return
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return
	}

	// read the HTML page as a PDF page
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(body.Bytes()))
	// enable this if the HTML file contains local references such as images, CSS, etc.
	page.EnableLocalFileAccess.Set(true)
	// add the page to your generator
	pdfg.AddPage(page)
	// manipulate page attributes as needed
	pdfg.MarginLeft.Set(10)
	pdfg.MarginRight.Set(10)
	pdfg.MarginTop.Set(15)
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	// magic
	err = pdfg.Create()
	if err != nil {
		return
	}

	var fileName string
	if invoice.Status == 2 || invoice.Status.String() == "Done" {
		fileName = "PP" + strconv.FormatInt(invoice.Id, 10) + "-" + invoice.StartDateInvoiceString() + ".pdf"
	} else {
		fileName = "PP" + strconv.FormatInt(invoice.Id, 10) + "-statement-" + invoice.StartDateInvoiceString() + ".pdf"
	}
	if !utility.IsWindow() {
		err = pdfg.WriteFile("/home/assyrian/go/selfserve/source/www/themes/muze/assets/invoice/" + fileName)
	} else {
		err = pdfg.WriteFile("D:/ProjectGo/PubPower/source/www/themes/muze/assets/invoice/" + fileName)
	}
	if err != nil {
		return
	}

	pathPDF = "/assets/invoice/" + fileName
	return
}

func setAddress(address, city, country string) (Address string) {
	if city != "" {
		city = ", " + city
	}
	if country != "" {
		country = ", " + country
	}

	Address = address + city + country
	return
}
