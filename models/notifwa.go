package models

type NotifPaymentWa struct {
	AggrNo string `json:"aggrNo"`
	Amount string `json:"amount"`
	WaNo   string `json:"waNo"`
}

type InstReqToMeta struct {
	MessagingProduct string   `json:"messaging_product"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Text             TextBody `json:"text"`
}

type TextBody struct {
	Body string `json:"body"`
}

type ResFrMeta struct {
	MessagingProduct string `json:"messaging_product"`
	Contacts         []struct {
		Input string `json:"input"`
		WaID  string `json:"wa_id"`
	} `json:"contacts"`
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
}
