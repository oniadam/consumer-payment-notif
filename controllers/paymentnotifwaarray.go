package controllers

import (
	"consumer-payment-notif/models"
	"consumer-payment-notif/services"
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func PaymentNotifWaArray(ch *amqp.Channel, logger *log.Logger) {
	msgs, err := ch.Consume(
		"paymentnotificationwa_queue",
		"",
		false, // manual ACK
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Println("Gagal mengambil pesan dari antrian:", err)

		return
	}

	data := []models.NotifPaymentWa{}

	for d := range msgs {

		func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Println("Panic recovered:", r)
					d.Nack(false, true)
				}
			}()

			logger.Println("Pesan yang diterima:", string(d.Body)+"queueName : paymentnotificationwa_queue")
			errs := json.Unmarshal(d.Body, &data)
			if errs != nil {
				logger.Println("Error unmarshaling JSON:", errs.Error())
				d.Nack(false, false) // supaya msk ke error queue
				// break
				return
			}
			datareceive := string(d.Body)
			resp, errSend := services.NotifPaymentWaArray(datareceive, data, logger)
			logger.Println("wefhwoehfowiehofwihefwewa", resp.ResponseCode)
			logger.Println("wefhwoehfowiehofwihefwewa", resp.ResponseMessage)
			logger.Println("wefhwoehfowiehofwihefwewa", resp.Errors)
			if errSend != nil {
				logger.Println(errSend.Error())
				// d.Nack(false, true)
				d.Nack(false, false) // supaya msk ke error queue
				// break
				return
				// continue
			}

			if resp.ResponseCode != "200" {
				logger.Printf("Query gagal, error=%s, code=%s, message=%s\n", resp.Errors, resp.ResponseCode, resp.ResponseMessage)

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
					logger.Println("Gagal requeue manual:", err.Error())
					d.Nack(false, true)
					return
				} else {
					d.Ack(false)
					return
				}

			}

			// gunakan kalo mau delay kirim pesan
			time.Sleep(3 * time.Second)

			d.Ack(false)
		}()

	}

}
