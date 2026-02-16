package services

import (
	"consumer-payment-notif/models"
	"consumer-payment-notif/repo"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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

	bodyMsg := ""
	if req.Flagreversal == "0" {
		// call sp utk get no reff
		datawano, _, _ := repo.GetReffNo(req.Refno)
		req.Refno = datawano.ReffNo

		bodyMsg = "Pembayaran angsuran Esta Dana Ventura Ibu " + req.CustomerName + " dengan no perjanjian " + req.AggrNo + " tanggal " + req.Senddtm + " sebesar " + fmt.Sprintf("%v", req.TotalPaid) + " (no transaksi " + req.Refno + " ) telah diterima. Untuk riwayat pembayaran dapat dilihat pada aplikasi MyEsta, Informasi lebih lanjut serta penyampaian pertanyaan atau keluhan, silahkan menghubungi Whatsapp di no nomor 081212xxxx"
	} else {
		bodyMsg = "Pembayaran angsuran Esta Dana Ventura Ibu " + req.CustomerName + " dengan no perjanjian " + req.AggrNo +
			" tanggal " + req.Senddtm + " sebesar " + fmt.Sprintf("%v", req.TotalPaid) + " (no transaksi " + req.Refno + " ) gagal tranksasi. Untuk riwayat pembayaran dapat dilihat pada aplikasi MyEsta, Informasi lebih lanjut serta penyampaian pertanyaan atau keluhan, silahkan menghubungi Whatsapp di no nomor 081212xxxx"
	}

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
		SetHeader("Authorization", os.Getenv("token_meta")).
		SetResult(&resfrmeta).
		// SetError(&resfrmetaerr).
		Post(os.Getenv("url_meta"))

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

	jsnResMeta, _ := json.Marshal(resfrmeta)
	log.Println("res fr meta", string(jsnResMeta))

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
