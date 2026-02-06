package controllers

import (
	"consumer-payment-notif/models"
	"consumer-payment-notif/services"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func PaymentNotifWa(ch *amqp.Channel, logger *log.Logger) {
	msgs, err := ch.Consume(
		"paymentnotificationwa_queue", // queue
		"",                            // consumer
		false,                         // auto-ack
		false,                         // exclusive
		false,                         // no-local
		false,                         // no-wait
		nil,                           // args
	)
	if err != nil {
		logger.Println("Gagal mengambil pesan dari antrian:", err)

		return
	}

	data := models.NotifPaymentWa{}

	for d := range msgs {
		logger.Println("Pesan yang diterima:", string(d.Body)+"queueName : paymentnotificationwa_queue")
		errs := json.Unmarshal(d.Body, &data)
		if errs != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}
		datareceive := string(d.Body)
		resp, errSend := services.NotifPaymentWa(data.AggrNo, datareceive, data)
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

		d.Ack(false)
	}

}
