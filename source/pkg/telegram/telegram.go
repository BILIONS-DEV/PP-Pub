package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/valyala/fasttemplate"
	"net/url"
	"source/pkg/ajax"
	"source/pkg/utility"
	"time"
)

type ResponseSendMessage struct {
	Ok     bool   `json:"ok"`
	Result Result `json:"result"`
}

type Result struct {
	MessageID int    `json:"message_id"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
}

type From struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type Chat struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

func SendMessageGroupPubPowerNotify(domainName string, user string, manager string) (err error) {
	if utility.IsDemo() || utility.IsDev() || utility.IsWindow() {
		return
	}
	link := "https://api.telegram.org/bot5095548297:AAGlcpxlhM9IuNZUcGPWpUoMAlrreO5m6Ng/sendMessage?chat_id=-1001441904341&parse_mode=html&text="

	template := "<b><u>New Domain</u></b>\n"
	template += "<i>-Domain</i>: {{domain}}\n"
	template += "<i>-User/Manager</i>: {{user_manager}}\n"
	template += "<i>-Link Similar</i>: {{link_similar}}\n"
	temp := fasttemplate.New(template, "{{", "}}")
	message := temp.ExecuteString(map[string]interface{}{
		"domain":       domainName,
		"user_manager": user + "/" + manager,
		"link_similar": "https://www.similarweb.com/website/" + domainName + "/#overview",
	})
	encodeMessage := url.QueryEscape(message)

	// makeRequest(c)
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	resp := ResponseSendMessage{}
	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &resp)
	})

	// Set error handler
	c.OnError(func(r *colly.Response, errHandle error) {
		if errHandle != nil {
			err = errHandle
		} else {
			if r.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("status code error: %d", r.StatusCode))
			}
		}
		if err != nil {
			return
		}
	})
	err = c.Visit(link + encodeMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Wait()
	return
}

func SendMessagePushLineGoogleError(notifyError string) (err error) {
	if utility.IsDemo() || utility.IsDev() || utility.IsWindow() {
		return
	}
	link := "https://api.telegram.org/bot5095548297:AAGlcpxlhM9IuNZUcGPWpUoMAlrreO5m6Ng/sendMessage?chat_id=-1001697924059&parse_mode=html&text="
	message := notifyError
	encodeMessage := url.QueryEscape(message)

	// makeRequest(c)
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	resp := ResponseSendMessage{}
	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &resp)
	})

	// Set error handler
	c.OnError(func(r *colly.Response, errHandle error) {
		if errHandle != nil {
			err = errHandle
		} else {
			if r.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("status code error: %d", r.StatusCode))
			}
		}
		if err != nil {
			return
		}
	})
	err = c.Visit(link + encodeMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Wait()
	return
}

func SendMessageNewVideo(notifyError string) (err error) {
	if utility.IsDemo() || utility.IsDev() || utility.IsWindow() {
		return
	}
	link := "https://api.telegram.org/bot5095548297:AAGlcpxlhM9IuNZUcGPWpUoMAlrreO5m6Ng/sendMessage?chat_id=-663055102&parse_mode=html&text="
	message := notifyError
	encodeMessage := url.QueryEscape(message)

	// makeRequest(c)
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	resp := ResponseSendMessage{}
	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &resp)
	})

	// Set error handler
	c.OnError(func(r *colly.Response, errHandle error) {
		if errHandle != nil {
			err = errHandle
		} else {
			if r.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("status code error: %d", r.StatusCode))
			}
		}
		if err != nil {
			return
		}
	})
	err = c.Visit(link + encodeMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Wait()
	return
}

func SendMessageTest(notifyError string) (err error) {
	link := "https://api.telegram.org/bot5095548297:AAGlcpxlhM9IuNZUcGPWpUoMAlrreO5m6Ng/sendMessage?chat_id=-1001616680965&parse_mode=html&text="
	message := notifyError
	encodeMessage := url.QueryEscape(message)

	// makeRequest(c)
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	resp := ResponseSendMessage{}
	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &resp)
	})

	// Set error handler
	c.OnError(func(r *colly.Response, errHandle error) {
		if errHandle != nil {
			err = errHandle
		} else {
			if r.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("status code error: %d", r.StatusCode))
			}
		}
		if err != nil {
			return
		}
	})
	err = c.Visit(link + encodeMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Wait()
	return
}

func SendErrorTuan(mesageError string, title ...string) (err error) {
	if utility.IsDemo() || utility.IsDev() || mesageError == "" {
		return
	}
	message := mesageError + "\n" + time.Now().Format("02-01-2006 15:04:05")
	if len(title) > 0 {
		message = title[0] + "\n" + message
	}
	link := "https://api.telegram.org/bot5404660010:AAHZPH7llQn_kCxw4jD9xHF94_L5va-X4Sg/sendMessage?chat_id=1073756657&parse_mode=html&text=" + url.QueryEscape(message)
	// makeRequest(c)
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	resp := ResponseSendMessage{}
	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &resp)
	})

	// Set error handler
	c.OnError(func(r *colly.Response, errHandle error) {
		if errHandle != nil {
			err = errHandle
		} else {
			if r.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("status code error: %d", r.StatusCode))
			}
		}
		if err != nil {
			return
		}
	})
	err = c.Visit(link)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Wait()
	return
}

func SendErrorsTuan(errs []ajax.Error, title ...string) (err error) {
	if utility.IsDemo() || utility.IsDev() || len(errs) == 0 {
		return
	}
	var mesageError string
	for _, value := range errs {
		mesageError = mesageError + "\n" + value.Id + " " + value.Message
	}
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return
	}
	message := mesageError + "\n" + time.Now().In(loc).Format("02-01-2006 15:04:05")
	if len(title) > 0 {
		message = title[0] + "\n" + message
	}
	link := "https://api.telegram.org/bot5404660010:AAHZPH7llQn_kCxw4jD9xHF94_L5va-X4Sg/sendMessage?chat_id=1073756657&parse_mode=html&text=" + url.QueryEscape(message)
	// makeRequest(c)
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	resp := ResponseSendMessage{}
	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &resp)
	})

	// Set error handler
	c.OnError(func(r *colly.Response, errHandle error) {
		if errHandle != nil {
			err = errHandle
		} else {
			if r.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("status code error: %d", r.StatusCode))
			}
		}
		if err != nil {
			return
		}
	})
	err = c.Visit(link)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Wait()
	return
}
