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

func NotifPaymentWaArray(datareceice string, req []models.NotifPaymentWa, logger *log.Logger) (res models.Respons, err error) {
	timeStr := time.Now().Format("2006-01-02 15:04:05")

	for _, v := range req {
		// validasi param request ke meta string kosong
		if v.AggrNo == "" || v.CustomerName == "" || v.Refno == "" || v.CustomerServiceNo == "" {
			res = models.Respons{
				ResponseCode:      "500",
				ResponseMessage:   "Parameter ada yang kosong",
				ResponseTimestamp: timeStr,
				Errors:            "",
				Data:              nil,
			}

			return res, err
		}

		if v.Createddtm == "" {
			v.Createddtm = timeStr
		}

		layout := "2006-01-02 15:04:05"

		t, err := time.Parse(layout, v.Senddtm)
		if err != nil {
			t, err = time.Parse(time.RFC3339Nano, v.Senddtm)
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

		totpaid := fmt.Sprintf("%v", v.TotalPaid)
		totpaid = utils.FormatRupiah(totpaid)

		// bodyMsg := ""
		reqtometa := models.NotifPaymentTemplateToMetaFbReq{}
		if v.Flagreversal == "0" {

			// datatempcode, _, _ := repo.GetTemplateCode("notif_payment")
			// fmt.Println(datatempcode)
			// req.Templatecode = datatempcode.TemplateCode

			reqtometa = models.NotifPaymentTemplateToMetaFbReq{
				MessagingProduct: "whatsapp",
				To:               v.WaNo,
				Type:             "template",
				Template: models.NotifPaymentTemplates{
					Name: v.Templatecode,
					LanguageS: models.LanguageS{
						Code: v.LanguageCode,
					},
					ComponentsS: []models.ComponentS{
						{
							Type: "body",
							Parameters: []models.ParameterS{
								{Type: "text", Text: v.CustomerName},
								{Type: "text", Text: v.AggrNo},
								{Type: "text", Text: outputsenddtmwa},
								{Type: "text", Text: totpaid},
								{Type: "text", Text: v.Refno},
								{Type: "text", Text: v.CustomerServiceNo},
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
				To:               v.WaNo,
				Type:             "template",
				Template: models.NotifPaymentTemplates{
					Name: v.Templatecode,
					LanguageS: models.LanguageS{
						Code: v.LanguageCode,
					},
					ComponentsS: []models.ComponentS{
						{
							Type: "body",
							Parameters: []models.ParameterS{
								{Type: "text", Text: v.CustomerName},
								{Type: "text", Text: v.AggrNo},
								{Type: "text", Text: outputsenddtmwa},
								{Type: "text", Text: totpaid},
								{Type: "text", Text: v.Refno},
								{Type: "text", Text: v.CustomerServiceNo},
							},
						},
					},
				},
			}

		}

		insPaymentNotifWa, errinsPaymentNotifWa := repo.InsertPaymentNotifWaRepo(outputsenddtmdb, v.Sendby, v.WaNo, v.Templatecode, v.AggrNo, v.CustomerName, fmt.Sprintf("%v", v.TotalPaid), v.TransactionSrc, v.Paymentmetodcode, v.Refno, v.Filepath, v.Flagreversal, v.Createdby, v.Createddtm)
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
		var urlkemata, tokenkemeta string

		// cek ukm/mikro
		if v.SenderHpNo == os.Getenv("senderHpNo1") || v.SenderHpNo == os.Getenv("senderHpNo2") {
			urlkemata = os.Getenv("url_meta_mikro")
			tokenkemeta = os.Getenv("token_meta_mikro")
		} else {
			urlkemata = os.Getenv("url_meta_ukm")
			tokenkemeta = os.Getenv("token_meta_ukm")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		restyClient := resty.New()
		resSendReminder, errSend := restyClient.R().
			SetContext(ctx).
			SetBody(&reqtometa).
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", tokenkemeta).
			SetResult(&resfrmeta).
			// SetError(&resfrmetaerr).
			Post(urlkemata)

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
		if v.Flagreversal == "0" {
			reqMyEsta = models.NotifMyestaReq{
				IdNo:             v.IdNo,
				Token:            "",
				Title:            "Pembayaran Berhasil",
				ShortDescription: "Pembayaran angsuran ESTA Dana Ventura telah kami terima",
				FullDescription:  "Pembayaran Diterima<br><br>Pembayaran angsuran <b>ESTA Dana Ventura</b> atas nama Ibu <b>" + v.CustomerName + "</b> dengan nomor perjanjian <b>" + v.AggrNo + "</b> pada tanggal <b>" + outputsenddtmwa + "</b> sebesar <b>Rp. " + totpaid + "</b> (Nomor Transaksi: <b>" + v.Refno + "</b>) telah kami terima.<br><br>Riwayat pembayaran dapat dilihat melalui menu <b>Riwayat Transaksi</b>, dan data pada <b>kartu piutang</b> akan diperbarui maksimal <b>H+1</b> setelah pembayaran berhasil dilakukan.<br><br>Untuk informasi lebih lanjut, pertanyaan, maupun keluhan, silakan menghubungi hotline kami melalui WhatsApp di <b>" + v.CustomerServiceNo + "</b> .<br><br>Terima kasih.",
				ResponseFcm:      "",
				Timestamp:        0,
				DeviceId:         "",
				CreatedBy:        v.Sendby,
			}
		} else {
			reqMyEsta = models.NotifMyestaReq{
				IdNo:             v.IdNo,
				Token:            "",
				Title:            "Pembayaran Gagal",
				ShortDescription: "Pembayaran angsuran ESTA Dana Ventura belum berhasil diproses/gagal",
				FullDescription:  "Pembayaran Gagal<br><br>Pembayaran angsuran <b>ESTA Dana Ventura</b> atas nama Ibu <b>" + v.CustomerName + "</b> dengan nomor perjanjian <b>" + v.AggrNo + "</b> pada tanggal <b>" + outputsenddtmwa + "</b> sebesar Rp. " + totpaid + " (Nomor Transaksi: <b>" + v.Refno + "</b>) belum berhasil diproses/gagal<br><br>Untuk informasi lebih lanjut, pertanyaan, maupun keluhan, silakan menghubungi hotline kami melalui WhatsApp di <b>" + v.CustomerServiceNo + "</b>.<br><br>Terima kasih.",
				ResponseFcm:      "",
				Timestamp:        0,
				DeviceId:         "",
				CreatedBy:        v.Sendby,
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
			logger.Println("resp error myesta", errSendMyEsta.Error())
			// c.JSON(http.StatusInternalServerError, res)
			// return res, errSendMyEsta
		}

		logger.Println("resp myesta", string(resSendMyEsta.Body()))

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
