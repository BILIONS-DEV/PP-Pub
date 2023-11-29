package mailjet

import (
	mailjets "github.com/mailjet/mailjet-apiv3-go"
	"source/pkg/utility"
)

const (
	apiKeyPublic  = "a03bbaab9cf3fdd0f57572a41158754f"
	apiKeyPrivate = "67c5c53eae43c0b266e6ccd5b80f05c6"
	// sub account Pubpower
	//API Key      3f006d3d25766809b0c22af446e580af
	//Secret Key   d50dcb8e95b5865a2c1a23e3e0d7f553
	FromEmail    = "noreply@pubpower.io"
	FromName     = "Noreply - Pubpower.io"
	FromEmailVli = "noreply@valueimpression.com"
	FromNameVli  = "Valueimpression"
)

var mailjetClient *mailjets.Client

type InfoMail struct {
	To          Email
	CC          []Email
	BCC         []Email
	Subject     string
	ContentText string
	ContentHtml string
}

type Email struct {
	Email string
	Name  string
}

func init() {
	mailjetClient = mailjets.NewMailjetClient(apiKeyPublic, apiKeyPrivate)
}

// custom[0] - Email, custom[1] - Name
func (this InfoMail) SendMail(custom ...string) error {
	if utility.IsWindow() {
		// Test local khong request
		return nil
	}
	var fromEmail, fromName string
	if len(custom) > 1 {
		fromEmail = custom[0]
		fromName = custom[1]
	} else {
		fromEmail = FromEmail
		fromName = FromName
	}
	inforMessage := mailjets.InfoMessagesV31{
		From: &mailjets.RecipientV31{
			Email: fromEmail,
			Name:  fromName,
		},
		To: &mailjets.RecipientsV31{
			mailjets.RecipientV31{
				Email: this.To.Email,
				Name:  this.To.Name,
			},
		},
		Subject:  this.Subject,
		TextPart: this.ContentText,
		HTMLPart: this.ContentHtml,
	}

	if len(this.CC) > 0 {
		sliceCc := mailjets.RecipientsV31{}
		for _, v := range this.CC {
			sliceCc = append(sliceCc, mailjets.RecipientV31{
				Email: v.Email,
				Name:  v.Name,
			})
		}
		inforMessage.Cc = &sliceCc
	}
	if len(this.BCC) > 0 {
		sliceBcc := mailjets.RecipientsV31{}
		for _, v := range this.CC {
			sliceBcc = append(sliceBcc, mailjets.RecipientV31{
				Email: v.Email,
				Name:  v.Name,
			})
		}
		inforMessage.Bcc = &sliceBcc
	}

	messagesInfo := []mailjets.InfoMessagesV31{inforMessage}
	messages := mailjets.MessagesV31{Info: messagesInfo}
	_, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		return err
	}

	return nil
}
