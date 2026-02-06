package services

import (
	"consumer-payment-notif/models"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

func NotifPaymentWa(aggrno string, datareceice string, req models.NotifPaymentWa) (res models.Respons, err error) {

	bodyMsg := "Pembayaran Anda dengan No Kontrak " + req.AggrNo + " telah terbayarkan senilai Rp. " + req.Amount

	reqtometa := models.InstReqToMeta{
		MessagingProduct: "whatsapp",
		To:               req.WaNo,
		Type:             "text",
		Text: models.TextBody{
			Body: bodyMsg,
		},
	}

	jsnLogMeta, _ := json.Marshal(reqtometa)
	log.Println("req to meta", string(jsnLogMeta))

	resfrmeta := models.ResFrMeta{}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	restyClient := resty.New()
	resSendReminder, errSend := restyClient.R().
		SetContext(ctx).
		SetBody(&reqtometa).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer EAATl4Ljd4UoBP3ZAeNMSRibZAwh1ucXHeu0hZBnqU3m9Gnv6UvymP7ZBqPFxYCMQHnzURqRfPcfhji2ZAKga65q1wjLS8B0pPkwFo4DZC5iVqnDeXH6ippO4XLAcEJrIg7S7loC1TTZCjLpo7jd8VZA0Kcxi41FlxfV7NLNvjS5XjpqlEqm1G8hg2GHvUYvmWdxy5wZDZD").
		SetResult(&resfrmeta).
		// SetError(&resfrmetaerr).
		Post("https://graph.facebook.com/v22.0/860450543810738/messages")

	if errSend != nil {
		res = models.Respons{
			ResponseCode:      "500",
			ResponseMessage:   "Internal Server Error",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            errSend.Error(),
			Data:              nil,
		}
		// c.JSON(http.StatusInternalServerError, res)
		return res, errSend
	}

	if resSendReminder.StatusCode() != 200 {
		res = models.Respons{
			ResponseCode:      "500",
			ResponseMessage:   "Internal Server Error",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            "resSendReminder.StatusCode() != 200",
			Data:              nil,
		}
		// c.JSON(http.StatusInternalServerError, res)
		return res, err
	}

	jsnResMeta, _ := json.Marshal(resfrmeta)
	log.Println("res fr meta", string(jsnResMeta))

	res = models.Respons{
		ResponseCode:      "200",
		ResponseMessage:   "success",
		ResponseTimestamp: "",
		Errors:            "",
		Data:              nil,
	}

	return res, nil
}
