package controllers

import (
	"consumer-payment-notif/models"
	"consumer-payment-notif/repo"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func PaymentNotif(ch *amqp.Channel, logger *log.Logger) {
	msgs, err := ch.Consume(
		"paymentnotification_queue", // queue
		"",                          // consumer
		false,                       // auto-ack
		false,                       // exclusive
		false,                       // no-local
		false,                       // no-wait
		nil,                         // args
	)
	if err != nil {
		logger.Println("Gagal mengambil pesan dari antrian:", err)

		return
	}

	data := models.NotifPayment{}

	for d := range msgs {
		logger.Println("Pesan yang diterima:", string(d.Body))
		errs := json.Unmarshal(d.Body, &data)
		if errs != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}
		datareceive := string(d.Body)
		resp, errSend := NotifPayment(data.AggrNo, datareceive, data)
		fmt.Println("wefhwoehfowiehofwihefwe", resp.ResponseCode)
		if errSend != nil {
			logger.Println(errSend.Error())
			d.Nack(false, true)
			// break
			return
			// continue
		}

		if resp.ResponseCode != "200" {
			logger.Printf("Query gagal, code=%s, message=%s\n", resp.ResponseCode, resp.ResponseMessage)

			// d.Nack(false, true)
			// return
			err := ch.Publish(
				"",           // exchange default
				d.RoutingKey, // sama dengan queue
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        d.Body,
				},
			)
			if err != nil {
				logger.Println("Gagal requeue manual:", err)
				d.Nack(false, true)
				return
			} else {
				d.Ack(false)
				return
			}

		}

		// gunakan kalo mau delay kirim pesan
		time.Sleep(15 * time.Second)
		logger.Println("Berhasil insert ke DB untuk AggrNo:", data.AggrNo)

		d.Ack(false)
	}

}

func NotifPayment(aggrno string, datareceice string, data models.NotifPayment) (res models.Respons, err error) {
	// var c *gin.Context

	// titlenotif := "Pembayaran Anda dengan No Kontrak " + data.AggrNo + " telah terbayarkan senilai Rp. " + data.Amount
	// descnotif := "Pembayaran Anda dengan No Kontrak " + data.AggrNo + " telah terbayarkan senilai Rp. " + data.Amount
	// fulldescnotif := "Pembayaran Anda dengan No Kontrak " + data.AggrNo + " telah terbayarkan senilai Rp. " + data.Amount

	// resp, errInsert := repo.InsertNotificationRepo(titlenotif, descnotif, fulldescnotif)
	resp, errInsert := repo.InsertTEstRabit(aggrno, data.Amount)
	if errInsert != nil {
		res = models.Respons{
			ResponseCode:    "500",
			ResponseMessage: "error query" + errInsert.Error(),
		}
		// c.JSON(500, res)
		return res, errInsert
	}
	fmt.Println("rererwerwerwerw", resp.ResponseCode)
	if resp.ResponseCode != "200" {
		res = models.Respons{
			ResponseCode:    resp.ResponseCode,
			ResponseMessage: resp.ResponseMessage,
		}
		// c.JSON(500, res)
		return res, nil
	}

	res = models.Respons{

		ResponseCode:    "200",
		ResponseMessage: "success",
	}

	return res, nil
}

func FlushLoggerIfFile(logger *log.Logger) {
	if logger == nil {
		return
	}
	w := logger.Writer()
	if f, ok := w.(*os.File); ok {
		_ = f.Sync() // ignore error, tapi bisa log kalau mau
	}
}
