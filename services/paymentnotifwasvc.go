package services

import (
	"consumer-payment-notif/models"
	"consumer-payment-notif/repo"
	"consumer-payment-notif/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

func NotifPaymentWa(aggrno string, datareceice string, req models.NotifPaymentWa, logger *log.Logger) (res models.Respons, err error) {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	// var datawano models.GetWaNoRes
	// if req.WaNo == "null" || req.WaNo == "" {
	// 	datawano, _, _ = repo.GetWaNo(req.AggrNo)
	// } else {
	// 	datawano, _, _ = repo.GetWaNo(req.AggrNo)
	// 	datawano.WaNo.String = req.WaNo
	// }

	if req.Createddtm == "" {
		req.Createddtm = timeStr
	}

	layout := "2006-01-02 15:04:05"

	t, err := time.Parse(layout, req.Senddtm)
	if err != nil {
		t, err = time.Parse(time.RFC3339Nano, req.Senddtm)
		if err != nil {
			res = models.Respons{
				ResponseCode:      "500",
				ResponseMessage:   "Error Parse Time ",
				ResponseTimestamp: timeStr,
				Errors:            err.Error() + "Parse",
				Data:              nil,
			}

			return res, err
		}

	}

	bulanIndo := map[time.Month]string{
		time.January:   "Januari",
		time.February:  "Februari",
		time.March:     "Maret",
		time.April:     "April",
		time.May:       "Mei",
		time.June:      "Juni",
		time.July:      "Juli",
		time.August:    "Agustus",
		time.September: "September",
		time.October:   "Oktober",
		time.November:  "November",
		time.December:  "Desember",
	}

	outputsenddtmwa := fmt.Sprintf("%d %s %d",
		t.Day(),
		bulanIndo[t.Month()],
		t.Year(),
	)

	outputsenddtmdb := t.Format("2006-01-02 15:04:05")

	totpaid := fmt.Sprintf("%v", req.TotalPaid)
	totpaid = utils.FormatRupiah(totpaid)

	// bodyMsg := ""
	reqtometa := models.NotifPaymentTemplateToMetaFbReq{}
	if req.Flagreversal == "0" {

		// datatempcode, _, _ := repo.GetTemplateCode("notif_payment")
		// fmt.Println(datatempcode)
		// req.Templatecode = datatempcode.TemplateCode
		reqtometa = models.NotifPaymentTemplateToMetaFbReq{
			MessagingProduct: "whatsapp",
			To:               req.WaNo,
			Type:             "template",
			Template: models.NotifPaymentTemplates{
				Name: req.Templatecode,
				LanguageS: models.LanguageS{
					Code: req.LanguageCode,
				},
				ComponentsS: []models.ComponentS{
					{
						Type: "body",
						Parameters: []models.ParameterS{
							{Type: "text", Text: req.CustomerName},
							{Type: "text", Text: req.AggrNo},
							{Type: "text", Text: outputsenddtmwa},
							{Type: "text", Text: totpaid},
							{Type: "text", Text: req.Refno},
							{Type: "text", Text: req.CustomerServiceNo},
						},
					},
				},
			},
		}

	} else {

		// datatempcode, _, _ := repo.GetTemplateCode("notif_reversal_payment")
		// fmt.Println(datatempcode)
		// req.Templatecode = datatempcode.TemplateCode
		reqtometa = models.NotifPaymentTemplateToMetaFbReq{
			MessagingProduct: "whatsapp",
			To:               req.WaNo,
			Type:             "template",
			Template: models.NotifPaymentTemplates{
				Name: req.Templatecode,
				LanguageS: models.LanguageS{
					Code: req.LanguageCode,
				},
				ComponentsS: []models.ComponentS{
					{
						Type: "body",
						Parameters: []models.ParameterS{
							{Type: "text", Text: req.CustomerName},
							{Type: "text", Text: req.AggrNo},
							{Type: "text", Text: outputsenddtmwa},
							{Type: "text", Text: totpaid},
							{Type: "text", Text: req.Refno},
							{Type: "text", Text: req.CustomerServiceNo},
						},
					},
				},
			},
		}

	}

	insPaymentNotifWa, errinsPaymentNotifWa := repo.InsertPaymentNotifWaRepo(outputsenddtmdb, req.Sendby, req.WaNo, req.Templatecode, req.AggrNo, req.CustomerName, fmt.Sprintf("%v", req.TotalPaid), req.TransactionSrc, req.Paymentmetodcode, req.Refno, req.Filepath, req.Flagreversal, req.Createdby, req.Createddtm)
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

	jsnLogMeta, _ := json.Marshal(reqtometa)
	logger.Println("req to meta", string(jsnLogMeta))

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
	logger.Println("res fr meta", string(jsnResMeta))

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

	// get token fcm user

	// request ke notif myesta mobile
	reqMyEsta := models.NotifMyestaReq{}
	if req.Flagreversal == "0" {
		reqMyEsta = models.NotifMyestaReq{
			IdNo:             req.IdNo,
			Token:            "",
			Title:            "Pembayaran Berhasil",
			ShortDescription: "Pembayaran Berhasil",
			FullDescription:  "Terimakasih telah melakukan pembayaran sebesar " + totpaid + " untuk no kontrak " + req.AggrNo,
			ResponseFcm:      "",
			Timestamp:        0,
			DeviceId:         "",
			CreatedBy:        req.Sendby,
		}
	} else {
		reqMyEsta = models.NotifMyestaReq{
			IdNo:             req.IdNo,
			Token:            "",
			Title:            "Pembayaran Gagal",
			ShortDescription: "Pembayaran Gagal",
			FullDescription:  "Gagal melakukan transaksi sebesar " + totpaid + " untuk no kontrak " + req.AggrNo,
			ResponseFcm:      "",
			Timestamp:        0,
			DeviceId:         "",
			CreatedBy:        req.Sendby,
		}
	}

	jsnReqMyesta, _ := json.Marshal(reqMyEsta)
	logger.Println("req to myesta", string(jsnReqMyesta))

	resSendMyEsta, errSendMyEsta := restyClient.R().
		SetContext(ctx).
		SetBody(&reqMyEsta).
		SetHeader("Content-Type", "application/json").
		SetResult(&res).
		Post(os.Getenv("url_notif_myesta"))

	if errSendMyEsta != nil {
		res = models.Respons{
			ResponseCode:      "500",
			ResponseMessage:   "Internal Server Error",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            errSendMyEsta.Error(),
			Data:              nil,
		}
		// c.JSON(http.StatusInternalServerError, res)
		return res, errSendMyEsta
	}

	logger.Println("resp myesta", string(resSendMyEsta.Body()))

	res = models.Respons{
		ResponseCode:      "200",
		ResponseMessage:   "success",
		ResponseTimestamp: "",
		Errors:            "",
		Data:              nil,
	}

	return res, nil
}
