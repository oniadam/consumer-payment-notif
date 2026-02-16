package repo

import (
	"consumer-payment-notif/constanta"
	"consumer-payment-notif/db"
	"consumer-payment-notif/models"
	"context"
	"database/sql"
	"time"
)

func GetReffNo(trxno string) (res models.GetReffNoRes, resError models.Respons, err error) {

	db, errcon := db.GetsSQLsrvDB3()

	if errcon != nil {
		resError = models.Respons{
			ResponseCode:      constanta.KODE_ERROR_KONEKSI_DB,
			ResponseMessage:   constanta.DESC_ERROR_KONEKSI_DB,
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            constanta.DESC_ERROR_KONEKSI_DB,
			Data:              nil,
		}

		return res, resError, errcon
	}

	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), constanta.TIMEOUT*time.Second)
	defer cancel()

	// Begin a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resError = models.Respons{
			ResponseCode:      constanta.DESC_ERROR_TIMEOUT,
			ResponseMessage:   constanta.DESC_ERROR_TIMEOUT,
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            constanta.DESC_ERROR_TIMEOUT,
			Data:              nil,
		}
		return res, resError, err
	}

	// Rollback the transaction on function exit
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, "exec spe_get_receipt_no  @pTransactionNo=?", trxno).Scan(&res.ReffNo)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, resError, err
		}

		if ctx.Err() == context.DeadlineExceeded {
			resError = models.Respons{
				ResponseCode:      constanta.KODE_ERROR_TIMEOUT,
				ResponseMessage:   constanta.DESC_ERROR_TIMEOUT,
				ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
				Errors:            constanta.DESC_ERROR_TIMEOUT,
				Data:              nil,
			}
			return res, resError, err
		}

		resError = models.Respons{
			ResponseCode:      constanta.KODE_ERROR_SP,
			ResponseMessage:   constanta.DESC_ERROR_SP,
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            constanta.DESC_ERROR_SP,
			Data:              nil,
		}
		return res, resError, err
	}

	tx.Commit()

	return res, resError, nil
}
