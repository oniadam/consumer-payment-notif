package repo

import (
	"consumer-payment-notif/db"
	"consumer-payment-notif/models"
	"context"
	"database/sql"
	"time"
)

func InsertPaymentNotifWaRepo(senddtm, sendby, wano, templatecode, aggrno, custname, totalpaid, transactionSrc, paymentmetodcode, refno, filepath, flagreversal, createdby, createddtm string) (res models.Respons, err error) {

	db, errcon := db.GetsSQLsrvDB3()

	if errcon != nil {
		res = models.Respons{
			ResponseCode:      "500",
			ResponseMessage:   "Error Koneksi DB",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            "",
			Data:              nil,
		}
		return res, errcon
	}

	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Begin a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		res = models.Respons{
			ResponseCode:      "408",
			ResponseMessage:   "ERROR_TIMEOUT",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            "",
			Data:              nil,
		}
		return res, err
	}

	// Rollback the transaction on function exit
	defer tx.Rollback()
	// layout := "2006-01-02 15:04:05"

	// // parse string ke time.Time
	// t, err := time.Parse(layout, senddtm)
	// // if err != nil {
	// // 	fmt.Println("Error:", err)
	// // 	return
	// // }

	// cdtm, _ := time.Parse(layout, createdtm)

	err = tx.QueryRowContext(ctx, "exec [spe_wa_send_notif_payment_insert] @pSendDtm=?, @pSendBy =?, @pWaNo=?, @pTemplateCode=?, @pAggrNo=?, @pCustomerFullname=?, @pTotalPaid=?, @pTransactionSource=?,  @pPaymentMethodCode=?, @pRefNo=?,@pFilePath=?, @pFlagReversal=?, @pCreatedBy=?, @pCreatedDtm=?", senddtm, sendby, wano, templatecode, aggrno, custname, totalpaid, transactionSrc, paymentmetodcode, refno, filepath, flagreversal, createdby, createddtm).Scan(&res.ResponseCode, &res.ResponseMessage, &res.Errors, &res.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		if ctx.Err() == context.DeadlineExceeded {
			res = models.Respons{
				ResponseCode:      "408",
				ResponseMessage:   "ERROR_TIMEOUT QUERY",
				ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
				Errors:            "",
				Data:              nil,
			}
			return res, err
		}

		res = models.Respons{
			ResponseCode:      "500",
			ResponseMessage:   "Terjadi Kendala System (1001)",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            "",
			Data:              nil,
		}
		return res, err
	}

	tx.Commit()

	return res, nil
}

func UpdatePaymentNotifWaRepo(sendreminderno, flag, chatid, updateby, updatedtm, rescode, jsonres string) (res models.Respons, err error) {

	db, errcon := db.GetsSQLsrvDB3()

	if errcon != nil {
		res = models.Respons{
			ResponseCode:      "500",
			ResponseMessage:   "Error Koneksi DB",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            "",
			Data:              nil,
		}
		return res, errcon
	}

	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Begin a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		res = models.Respons{
			ResponseCode:      "408",
			ResponseMessage:   "ERROR_TIMEOUT",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            "",
			Data:              nil,
		}
		return res, err
	}

	// Rollback the transaction on function exit
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, "exec [spe_wa_notif_payment_update] @pSendNotifPaymentNo=?, @pFlagFailed=?, @pRefNo=?, @pUpdatedBy=?, @pUpdatedDtm=?, @pResponseCode =?, @pJsonRes=?", sendreminderno, flag, chatid, updateby, updatedtm, rescode, jsonres).Scan(&res.ResponseCode, &res.ResponseMessage, &res.Errors, &res.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		if ctx.Err() == context.DeadlineExceeded {
			res = models.Respons{
				ResponseCode:      "408",
				ResponseMessage:   "ERROR_TIMEOUT QUERY",
				ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
				Errors:            "",
				Data:              nil,
			}
			return res, err
		}

		res = models.Respons{
			ResponseCode:      "500",
			ResponseMessage:   "Terjadi Kendala System (1001)",
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            "",
			Data:              nil,
		}
		return res, err
	}

	tx.Commit()

	return res, nil
}
