package main

import (
	"consumer-payment-notif/controllers"
	"consumer-payment-notif/utils"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func main() {
	// ========================================================================================
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err.Error())
	}

	// router.Use(utils.RequestLogger())

	// log := utils.InitLogger()

	// logger := utils.SetupLogger()
	utils.RotateLogger()     // init pertama
	utils.StartLogRotation() // auto rotate tiap 00:00

	logger := utils.GetLogger()

	// router := gin.Default()
	for {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			logger.Println("Gagal terhubung ke RabbitMQ:", err)
			// return
			time.Sleep(5 * time.Second)
			continue
		}
		logger.Println("Terhubung ke RabbitMQ")

		// chNotif, err := conn.Channel()
		// if err != nil {
		// 	logger.Println("Gagal membuka channel:", err)
		// 	conn.Close()
		// 	time.Sleep(5 * time.Second)
		// 	continue
		// }

		// err = chNotif.Qos(1, 0, false)
		// if err != nil {
		// 	logger.Println("Gagal set QoS:", err)
		// }

		// // Jalankan consumer
		// go controllers.PaymentNotif(chNotif, logger)

		// kalo mau buat bnyk consumer, buat lagi sperti chNotif
		// ================== Channel Consumer 2 ==================
		chNotifWa, err := conn.Channel()
		if err != nil {
			// logger.Println("Gagal membuka channel PaymentNotifWa:", err)
			// chNotifWa.Close()
			// conn.Close()
			// time.Sleep(5 * time.Second)
			// continue
			conn.Close()
			return
		}
		err = chNotifWa.Qos(1, 0, false)
		if err != nil {
			// logger.Println("Gagal set QoS PaymentNotifWa:", err)
			chNotifWa.Close()
			conn.Close()
			return
		}
		go controllers.PaymentNotifWaArray(chNotifWa, logger)

		// Tunggu sampai koneksi error
		errChan := make(chan *amqp.Error)
		conn.NotifyClose(errChan)

		// Tunggu error dari RabbitMQ (blocking)
		errConn := <-errChan
		if errConn != nil {
			logger.Println("⚠️ Koneksi RabbitMQ terputus:", errConn)
		}

		// Tutup koneksi dan tunggu 5 detik sebelum reconnect
		// chNotif.Close() // tutup lgi yg lain klo ada 3 konsmer
		chNotifWa.Close()
		conn.Close()
		// time.Sleep(5 * time.Second)
	}

	// log.Println(http.ListenAndServe(":8912", router))

	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// if err != nil {
	// 	log.Println("Gagal terhubung ke RabbitMQ:", err)
	// }
	// defer conn.Close()

	// // res := models.Respons{}

	// // Membuka channel
	// ch, err := conn.Channel()
	// if err != nil {
	// 	log.Println("Gagal membuka channel:", err)
	// }
	// defer ch.Close()

	// go controllers.PaymentNotif(ch, logger)
	// == tanpa reconnect

	// // Mengambil pesan dari antrian
	// msgs, err := ch.Consume(
	// 	"paymentnotification_queue", // queue
	// 	"",                          // consumer
	// 	false,                       // auto-ack
	// 	false,                       // exclusive
	// 	false,                       // no-local
	// 	false,                       // no-wait
	// 	nil,                         // args
	// )
	// if err != nil {
	// 	log.Println("Gagal mengambil pesan dari antrian:", err)
	// }

	// forever := make(chan bool)
	// data := models.NotifPayment{}
	// go func() {
	// 	for d := range msgs {
	// 		log.Println("Pesan yang diterima:", string(d.Body))
	// 		errs := json.Unmarshal(d.Body, &data)
	// 		if errs != nil {
	// 			fmt.Println("Error unmarshaling JSON:", err)
	// 			return
	// 		}
	// 		datareceive := string(d.Body)
	// 		errSend := controllers.SendNotif(data.AggrNo, datareceive, data)
	// 		if errSend != nil {
	// 			log.Println(errSend.Error())
	// 			return
	// 		}

	// 		// gunakan jika mau delay kirim pesan
	// 		time.Sleep(1 * time.Minute)
	// 		d.Ack(false)
	// 	}
	// }()

	// fmt.Println("Successfully Connected to our RabbitMQ Instance")
	// fmt.Println(" [*] - Waiting for messages")
	// <-forever

	// // c.String(http.StatusOK, "Menerima pesan dari RabbitMQ")

	// fmt.Println("Running...")

	// ===========================================================================================

}
