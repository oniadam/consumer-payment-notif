package repo

import (
	"consumer-payment-notif/constanta"
	"consumer-payment-notif/db"
	"consumer-payment-notif/models"
	"context"
	"database/sql"
	"time"
)

func InsertNotificationRepo(titlenotif, descnotif, fulldescnotif string) (resError models.Respons, err error) {
	// resError := models.Respons{}

	db, errcon := db.GetsSQLsrvDB()

	if errcon != nil {
		resError = models.Respons{
			ResponseCode:      constanta.KODE_ERROR_KONEKSI_DB,
			ResponseMessage:   constanta.DESC_ERROR_KONEKSI_DB,
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            constanta.DESC_ERROR_KONEKSI_DB,
			Data:              nil,
		}
		return resError, errcon
	}

	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), constanta.TIMEOUT*time.Second)
	defer cancel()

	// Begin a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resError = models.Respons{
			ResponseCode:      constanta.KODE_ERROR_TIMEOUT,
			ResponseMessage:   constanta.DESC_ERROR_TIMEOUT,
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            constanta.DESC_ERROR_TIMEOUT,
			Data:              nil,
		}
		return resError, err
	}

	// Rollback the transaction on function exit
	defer tx.Rollback()

	now := time.Now()

	// Tambah 1 hari
	tomorrow := now.AddDate(0, 0, 1)
	tomorrows := tomorrow.Format("2006-01-02 15:05:05")

	err = tx.QueryRowContext(ctx, "exec [spa_insert_send_notif] @pTitleNotifHp=?, @pDescNotifHp=?, @pTitleNotif=?, @pDescNotif=?,@pFullDescNotif=?,@pRefCode=?, @pExpDate=?, @textButton=?, @redirectTo=?, @descRedirect=?, @redirectType=?, @imageBase64=?, @pCreatedBy=?, @pCreatedDate=?", "", "", titlenotif, descnotif, fulldescnotif, "", tomorrows, "", "", "", "", "", "", tomorrows).Scan(&resError.ResponseCode, &resError.ResponseMessage, &resError.Errors, &resError.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return resError, nil
		}
		if ctx.Err() == context.DeadlineExceeded {
			resError = models.Respons{
				ResponseCode:      constanta.KODE_ERROR_TIMEOUT,
				ResponseMessage:   constanta.DESC_ERROR_TIMEOUT,
				ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
				Errors:            constanta.DESC_ERROR_TIMEOUT,
				Data:              nil,
			}
			return resError, err
		}

		resError = models.Respons{
			ResponseCode:      constanta.KODE_ERROR_SP,
			ResponseMessage:   constanta.DESC_ERROR_SP,
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            constanta.DESC_ERROR_SP,
			Data:              nil,
		}
		return resError, err
	}

	tx.Commit()

	return resError, nil
}

func InsertTEstRabit(aggrno, amount string) (resError models.Respons, err error) {
	// resError := models.Respons{}

	db, errcon := db.GetsSQLsrvDB2()

	if errcon != nil {
		resError = models.Respons{
			ResponseCode:      constanta.KODE_ERROR_KONEKSI_DB,
			ResponseMessage:   constanta.DESC_ERROR_KONEKSI_DB,
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            constanta.DESC_ERROR_KONEKSI_DB,
			Data:              nil,
		}
		return resError, errcon
	}

	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), constanta.TIMEOUT*time.Second)
	defer cancel()

	// Begin a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		resError = models.Respons{
			ResponseCode:      constanta.KODE_ERROR_TIMEOUT,
			ResponseMessage:   constanta.DESC_ERROR_TIMEOUT,
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            constanta.DESC_ERROR_TIMEOUT,
			Data:              nil,
		}
		return resError, err
	}

	// Rollback the transaction on function exit
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, "exec [spa_insert_rabbit] @pAggrNo=?, @pAmount=?", aggrno, amount).Scan(&resError.ResponseCode, &resError.ResponseMessage, &resError.Errors, &resError.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return resError, nil
		}
		if ctx.Err() == context.DeadlineExceeded {
			resError = models.Respons{
				ResponseCode:      constanta.KODE_ERROR_TIMEOUT,
				ResponseMessage:   constanta.DESC_ERROR_TIMEOUT,
				ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
				Errors:            constanta.DESC_ERROR_TIMEOUT,
				Data:              nil,
			}
			return resError, err
		}

		resError = models.Respons{
			ResponseCode:      constanta.KODE_ERROR_SP,
			ResponseMessage:   constanta.DESC_ERROR_SP,
			ResponseTimestamp: time.Now().Format("2006-01-02 15:04:05"),
			Errors:            constanta.DESC_ERROR_SP,
			Data:              nil,
		}
		return resError, err
	}

	tx.Commit()

	return resError, nil
}
