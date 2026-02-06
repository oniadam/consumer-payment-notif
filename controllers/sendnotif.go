package controllers

import (
	"consumer-payment-notif/models"
	"consumer-payment-notif/repo"

	"github.com/gin-gonic/gin"
)

func SendNotif(aggrno string, datareceice string, data models.NotifPayment) (err error) {
	var c *gin.Context

	titlenotif := "Pembayaran Anda dengan No Kontrak " + data.AggrNo + " telah terbayarkan senilai Rp. " + data.Amount
	descnotif := "Pembayaran Anda dengan No Kontrak " + data.AggrNo + " telah terbayarkan senilai Rp. " + data.Amount
	fulldescnotif := "Pembayaran Anda dengan No Kontrak " + data.AggrNo + " telah terbayarkan senilai Rp. " + data.Amount

	_, errInsert := repo.InsertNotificationRepo(titlenotif, descnotif, fulldescnotif)
	if errInsert != nil {
		c.JSON(500, err)
		return err
	}

	return nil
}
