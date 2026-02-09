package services

import (
	"consumer-payment-notif/models"
	"consumer-payment-notif/repo"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

func NotifPaymentWa(aggrno string, datareceice string, req models.NotifPaymentWa) (res models.Respons, err error) {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	totpaid := fmt.Sprintf("%v", req.TotalPaid)
	insPaymentNotifWa, errinsPaymentNotifWa := repo.InsertPaymentNotifWaRepo(req.Senddtm, req.Sendby, req.WaNo, req.Templatecode, req.AggrNo, req.CustomerName, totpaid, req.TransactionSrc, req.Paymentmetodcode, req.Refno, req.Filepath, req.Flagreversal, req.Createdby, req.Createddtm)
	if errinsPaymentNotifWa != nil {
		res = models.Respons{
			ResponseCode:      insPaymentNotifWa.ResponseCode,
			ResponseMessage:   insPaymentNotifWa.ResponseMessage,
			ResponseTimestamp: timeStr,
			Errors:            errinsPaymentNotifWa.Error() + "Sp Insert Reminder",
			Data:              nil,
		}

		// c.JSON(500, res)

		return res, errinsPaymentNotifWa
	}

	bodyMsg := "Hallo " + req.CustomerName + ", Pembayaran Anda dengan No Kontrak " + req.AggrNo + " telah terbayarkan senilai Rp. " + totpaid

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

	// save req data meta ke db

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

		jsnResMeta, _ := json.Marshal(resfrmeta)
		_, errUpdate := repo.UpdatePaymentNotifWaRepo(fmt.Sprintf("%v", insPaymentNotifWa.Data), "1", "", "system", time.Now().Format("2006-01-02 15:04:05.000"), strconv.Itoa(resSendReminder.StatusCode()), string(jsnResMeta))

		if errUpdate != nil {
			res = models.Respons{
				ResponseCode:      "500",
				ResponseMessage:   "Internal Server Error",
				ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
				Errors:            "Error Query Update " + errUpdate.Error(),
				Data:              nil,
			}
			// c.JSON(http.StatusInternalServerError, res)
			return res, err
		}

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

	// update req data meta ke db
	chatid := resfrmeta.Messages[0].ID

	_, errUpdate := repo.UpdatePaymentNotifWaRepo(fmt.Sprintf("%v", insPaymentNotifWa.Data), "0", chatid, "system", time.Now().Format("2006-01-02 15:04:05.000"), strconv.Itoa(resSendReminder.StatusCode()), string(jsnResMeta))

	if errUpdate != nil {
		res = models.Respons{
			ResponseCode:      "500",
			ResponseMessage:   "Internal Server Error",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            "Error Query Update " + errUpdate.Error(),
			Data:              nil,
		}
		// c.JSON(http.StatusInternalServerError, res)
		return res, err
	}

	res = models.Respons{
		ResponseCode:      "200",
		ResponseMessage:   "success",
		ResponseTimestamp: "",
		Errors:            "",
		Data:              nil,
	}

	return res, nil
}
